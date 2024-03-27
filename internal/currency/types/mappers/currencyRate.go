package mappers

import (
	"currencies/api"
	"currencies/internal/currency/types"

	"github.com/pkg/errors"
)

func DomainCurrencyRateToAPI(domainCurrencyRate types.CurrencyRate) (api.CurrencyRateResponse, error) {
	status, err := domainCurrencyRateStatusToAPI(domainCurrencyRate.Status)
	if err != nil {
		return api.CurrencyRateResponse{}, errors.Wrap(err, "mapping domain currency rate status to API")
	}

	return api.CurrencyRateResponse{
		Id:       domainCurrencyRate.Id,
		RateTime: domainCurrencyRate.RateTime,
		Status:   status,
		Rate:     domainCurrencyRate.Rate,
	}, nil
}

func domainCurrencyRateStatusToAPI(domainStatus types.CurrencyRateStatus) (api.CurrencyRateStatus, error) {
	switch domainStatus {
	case types.SUCCESS:
		return api.SUCCESS, nil
	case types.PENDING:
		return api.PENDING, nil
	case types.ERROR:
		return api.ERROR, nil
	default:
		return "", errors.Errorf("unknown currency rate status %s", domainStatus)
	}
}

func ApiCurrencyToDomain(apiCurrency api.Currency) (types.Currency, error) {
	switch apiCurrency {
	case api.EUR:
		return types.EUR, nil
	case api.USD:
		return types.USD, nil
	case api.MXN:
		return types.MXN, nil
	default:
		return "", errors.Errorf("unckown currency %s", apiCurrency)
	}
}
