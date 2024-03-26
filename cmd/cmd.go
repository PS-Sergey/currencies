package cmd

import (
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
	"go.uber.org/zap"
)

const (
	defaultConfigPath = "./config/config.yml"
)

func Run() {
	cfg := mustLoadConfig()
	logger := zap.Must(zap.NewProduction())

	db, err := cfg.DB.NewClient()
	if err != nil {
		logger.Fatal("Fail to get DB connection", zap.Error(err))
	}

	swagger, err := api.GetSwagger()
	if err != nil {
		logger.Fatal("Error loading swagger spec", zap.Error(err))
	}
	swagger.Servers = nil

	currencyRateRepository := repository.NewRateRepository(db)
	currencyRateUpdateChan := make(chan types.UpdateCurrencyRateMsg, 10)
	idGenerator := uuidGenerator.NewUUIDGenerator()
	currencyRateService := service.NewCurrencyRateService(currencyRateRepository, currencyRateUpdateChan, idGenerator)
	rateClient := cfg.ExchangeRateClient.NewClient()
	rateAdapter := exchangeRate.NewExchangeRateAdapter(rateClient)
	txManager := service.NewTransactionManager(db, logger)
	currencyRateUpdater := service.NewCurrencyRateUpdateService(logger, rateAdapter, txManager)
	currencyRateHandler := web.NewCurrencyRateHandler(currencyRateService, logger)

	r := mux.NewRouter()
	r.Use(middleware.OapiRequestValidator(swagger))
	handler := api.HandlerFromMux(currencyRateHandler, r)

	app := application.NewApp(cfg.Server, logger, currencyRateUpdateChan, currencyRateUpdater, handler)
	app.Run()
}

func mustLoadConfig() *config.Config {
	configPath := os.Getenv("CONFIG_PATH")
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