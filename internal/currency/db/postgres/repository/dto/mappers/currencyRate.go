package mappers

import (
	"currencies/internal/currency/db/postgres/repository/dto"
	"currencies/internal/currency/types"

	"github.com/pkg/errors"
)

func DomainCurrencyRateToDB(domainCurrencyRate types.CurrencyRate) (dto.CurrencyRate, error) {
	status, err := domainCurrencyRateStatusToDB(domainCurrencyRate.Status)
	if err != nil {
		return dto.CurrencyRate{}, errors.Wrap(err, "mapping currency rate status from domain to DB")
	}

	current, err := domainCurrencyToDB(domainCurrencyRate.Base)
	if err != nil {
		return dto.CurrencyRate{}, errors.Wrap(err, "mapping current currency from domain to DB")
	}

	target, err := domainCurrencyToDB(domainCurrencyRate.Target)
	if err != nil {
		return dto.CurrencyRate{}, errors.Wrap(err, "mapping target currency from domain to DB")
	}

	return dto.CurrencyRate{
		Id:       domainCurrencyRate.Id,
		RateTime: domainCurrencyRate.RateTime,
		Status:   status,
		Base:     current,
		Target:   target,
		Rate:     domainCurrencyRate.Rate,
	}, nil
}

func domainCurrencyRateStatusToDB(domainStatus types.CurrencyRateStatus) (dto.CurrencyRateStatus, error) {
	switch domainStatus {
	case types.PENDING:
		return dto.PENDING, nil
	case types.SUCCESS:
		return dto.SUCCESS, nil
	case types.ERROR:
		return dto.ERROR, nil
	default:
		return "", errors.Errorf("unknown currency rate status: %s", domainStatus)
	}
}

func domainCurrencyToDB(domainCurrency types.Currency) (dto.Currency, error) {
	switch domainCurrency {
	case types.EUR:
		return dto.EUR, nil
	case types.USD:
		return dto.USD, nil
	case types.MXN:
		return dto.MXN, nil
	default:
		return "", errors.Errorf("unknown currency: %s", domainCurrency)
	}
}

func DBCurrencyRateToDomain(dbCurrencyRate dto.CurrencyRate) (types.CurrencyRate, error) {
	status, err := dbCurrencyRateStatusToDomain(dbCurrencyRate.Status)
	if err != nil {
		return types.CurrencyRate{}, errors.Wrap(err, "mapping currency rate status from domain to DB")
	}

	current, err := dbCurrencyToDomain(dbCurrencyRate.Base)
	if err != nil {
		return types.CurrencyRate{}, errors.Wrap(err, "mapping current currency from domain to DB")
	}

	target, err := dbCurrencyToDomain(dbCurrencyRate.Target)
	if err != nil {
		return types.CurrencyRate{}, errors.Wrap(err, "mapping target currency from domain to DB")
	}

	return types.CurrencyRate{
		Id:       dbCurrencyRate.Id,
		RateTime: dbCurrencyRate.RateTime,
		Status:   status,
		Base:     current,
		Target:   target,
		Rate:     dbCurrencyRate.Rate,
	}, nil
}

func dbCurrencyRateStatusToDomain(dbStatus dto.CurrencyRateStatus) (types.CurrencyRateStatus, error) {
	switch dbStatus {
	case dto.PENDING:
		return types.PENDING, nil
	case dto.SUCCESS:
		return types.SUCCESS, nil
	case dto.ERROR:
		return types.ERROR, nil
	default:
		return "", errors.Errorf("unknown currency rate status: %s", dbStatus)
	}
}

func dbCurrencyToDomain(dbCurrency dto.Currency) (types.Currency, error) {
	switch dbCurrency {
	case dto.EUR:
		return types.EUR, nil
	case dto.USD:
		return types.USD, nil
	case dto.MXN:
		return types.MXN, nil
	default:
		return "", errors.Errorf("unknown currency: %s", dbCurrency)
	}
}
