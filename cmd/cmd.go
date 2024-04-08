package cmd

import (
	"context"
	"currencies/api"
	"currencies/cmd/application"
	"currencies/config"
	"currencies/internal/currency/clients/exchangeRate"
	"currencies/internal/currency/db/postgres/repository"
	"currencies/internal/currency/service"
	"currencies/internal/currency/transport/web"
	"currencies/internal/currency/types"
	"currencies/pkg/uuidGenerator"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/ilyakaznacheev/cleanenv"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.uber.org/zap"
)

const (
	configPathEnv     = "CONFIG_PATH"
	defaultConfigPath = "./config/config.yml"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := mustLoadConfig()
	logger := zap.Must(zap.NewProduction())

	err := cfg.Tracer.SetTracer(ctx)
	if err != nil {
		logger.Fatal("fail to set tracer", zap.Error(err))
	}

	db, err := cfg.DB.NewClient()
	if err != nil {
		logger.Fatal("fail to get DB connection", zap.Error(err))
	}
	defer db.Close()

	swagger, err := api.GetSwagger()
	if err != nil {
		logger.Fatal("error loading swagger spec", zap.Error(err))
	}
	swagger.Servers = nil

	currencyRateUpdateChan := make(chan types.UpdateCurrencyRateMsg, 10)
	defer close(currencyRateUpdateChan)

	currencyRateRepository := repository.NewRateRepository(db)
	idGenerator := uuidGenerator.NewUUIDGenerator()
	currencyRateService := service.NewCurrencyRateService(currencyRateRepository, currencyRateUpdateChan, idGenerator)
	rateClient := cfg.ExchangeRateClient.NewClient()
	rateAdapter := exchangeRate.NewExchangeRateAdapter(rateClient)
	txManager := service.NewTransactionManager(db, logger)
	currencyRateUpdater := service.NewCurrencyRateUpdateService(logger, rateAdapter, txManager)
	currencyRateHandler := web.NewCurrencyRateHandler(currencyRateService, logger)

	r := mux.NewRouter()
	r.Use(middleware.OapiRequestValidator(swagger))
	r.Use(otelmux.Middleware(cfg.Tracer.ServiceName))
	handler := api.HandlerFromMux(currencyRateHandler, r)

	app := application.NewApp(cfg.Server, logger, currencyRateUpdateChan, currencyRateUpdater, handler)
	app.Run(ctx)
}

func mustLoadConfig() *config.Config {
	configPath := os.Getenv(configPathEnv)
	if configPath == "" {
		configPath = defaultConfigPath
	}

	var c config.Config
	err := cleanenv.ReadConfig(configPath, &c)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &c
}
