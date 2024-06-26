package application

import (
	"context"
	"currencies/internal/currency"
	"currencies/internal/currency/types"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type App struct {
	cfg                    *Config
	logger                 *zap.Logger
	srv                    *http.Server
	handler                http.Handler
	currencyRateUpdateChan chan types.UpdateCurrencyRateMsg
	currencyRateUpdater    currency.RateUpdaterService
}

func NewApp(
	cfg *Config,
	logger *zap.Logger,
	currencyRateUpdateChan chan types.UpdateCurrencyRateMsg,
	currencyRateUpdater currency.RateUpdaterService,
	handler http.Handler,
) *App {
	return &App{
		cfg:                    cfg,
		logger:                 logger,
		currencyRateUpdateChan: currencyRateUpdateChan,
		currencyRateUpdater:    currencyRateUpdater,
		handler:                handler,
	}
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.withMiddleware()
	a.initServer()
	go a.runCurrencyRateUpdater(ctx)
	go func() {
		if err := a.srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				a.logger.Info("Server stopped")
				return
			}
			a.logger.Error("Server error", zap.Error(err))
		}
	}()
	a.logger.Info("server starting", zap.String("port", a.cfg.Port))

	a.gracefulShutdown(ctx)
}

func (a *App) runCurrencyRateUpdater(ctx context.Context) {
	for {
		select {
		case msg := <-a.currencyRateUpdateChan:
			a.currencyRateUpdater.UpdateCurrencyRate(ctx, msg)
		case <-ctx.Done():
			return
		}
	}
}

func (a *App) gracefulShutdown(ctx context.Context) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	a.srv.Shutdown(ctx)

	a.logger.Info("Shutting down")
	os.Exit(0)
}
