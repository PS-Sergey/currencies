package service

import (
	"context"
	"currencies/internal/currency"
	"currencies/internal/currency/types"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CurrencyRateService struct {
	currencyRateRepository currency.RateRepository
	updateCurrencyRateChan chan<- types.UpdateCurrencyRateMsg
	uuidGenerator          currency.UUIDGenerator
}

func NewCurrencyRateService(
	currencyRateRepository currency.RateRepository,
	updateCurrencyRateChan chan<- types.UpdateCurrencyRateMsg,
	uuidGenerator currency.UUIDGenerator,
) *CurrencyRateService {
	return &CurrencyRateService{
		currencyRateRepository: currencyRateRepository,
		updateCurrencyRateChan: updateCurrencyRateChan,
		uuidGenerator:          uuidGenerator,
	}
}

func (s *CurrencyRateService) NewCurrencyRate(ctx context.Context, base types.Currency, target types.Currency) (uuid.UUID, error) {
	currencyRateId, err := s.uuidGenerator.Generate()
	if err != nil {
		return uuid.UUID{}, errors.Wrap(err, "generate UUID for new currency rate")
	}

	newRate := types.CurrencyRate{
		Id:     currencyRateId,
		Status: types.PENDING,
		Base:   base,
		Target: target,
	}

	err = s.currencyRateRepository.SaveCurrencyRate(ctx, newRate)
	if err != nil {
		return uuid.UUID{}, errors.Wrap(err, "save new currency rate")
	}

	updateCurrencyRateMsg := types.UpdateCurrencyRateMsg{
		RequestId: currencyRateId,
		Base:      base,
		Target:    target,
	}

	s.updateCurrencyRateChan <- updateCurrencyRateMsg

	return currencyRateId, nil
}

func (s *CurrencyRateService) GetCurrencyRateById(ctx context.Context, currencyRateId uuid.UUID) (types.CurrencyRate, error) {
	rate, err := s.currencyRateRepository.GetCurrencyRateById(ctx, currencyRateId)
	if err != nil {
		return types.CurrencyRate{}, errors.Wrap(err, "get currency rate from repository")
	}

	return rate, nil
}

func (s *CurrencyRateService) GetLastCurrencyRate(ctx context.Context, base types.Currency, target types.Currency) (types.CurrencyRate, error) {
	rate, err := s.currencyRateRepository.GetLastCurrencyRate(ctx, base, target)
	if err != nil {
		return types.CurrencyRate{}, errors.Wrap(err, "get last currency rate from repository")
	}

	return rate, nil
}
