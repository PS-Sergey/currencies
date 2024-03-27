package service

import (
	"context"
	"currencies/internal/currency"
	"currencies/internal/currency/types"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type CurrencyRateUpdateService struct {
	logger             *zap.Logger
	exchangeRateClient currency.ExchangeRateClient
	txManager          currency.TransactionManager
}

func NewCurrencyRateUpdateService(
	logger *zap.Logger,
	exchangeRateClient currency.ExchangeRateClient,
	txManager currency.TransactionManager,
) *CurrencyRateUpdateService {
	return &CurrencyRateUpdateService{
		logger:             logger,
		exchangeRateClient: exchangeRateClient,
		txManager:          txManager,
	}
}

func (s *CurrencyRateUpdateService) UpdateCurrencyRate(ctx context.Context, msg types.UpdateCurrencyRateMsg) {
	err := s.txManager.BeginTx(ctx, nil, func(ctx context.Context, t currency.Transaction) error {
		return s.updateCurrencyRateTx(ctx, t, msg)
	})

	if err != nil {
		s.logger.Error("update currency rate", zap.String("id", msg.RequestId.String()), zap.Error(err))
	}
}

func (s *CurrencyRateUpdateService) updateCurrencyRateTx(ctx context.Context, t currency.Transaction, msg types.UpdateCurrencyRateMsg) error {
	txRateRepository := t.TxRateRepository()

	currencyRate, err := txRateRepository.GetCurrencyRateById(ctx, msg.RequestId)
	if err != nil {
		return errors.Wrap(err, "get currency rate from repository")
	}

	rate, err := s.exchangeRateClient.GetCurrencyRate(ctx, msg.Base, msg.Target)
	if err == nil {
		s.logger.Info("get currency rate from client", zap.Float32("rate:", rate))
		rateTime := time.Now()
		currencyRate.Status = types.SUCCESS
		currencyRate.Rate = &rate
		currencyRate.RateTime = &rateTime
	} else {
		s.logger.Error("get currency rate client request", zap.Error(err))
		currencyRate.Status = types.ERROR
	}

	err = txRateRepository.UpdateCurrencyRate(ctx, currencyRate)
	if err != nil {
		return errors.Wrap(err, "update currency rate")
	}

	return nil
}
