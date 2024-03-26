package types

import (
	"time"

	"github.com/google/uuid"
)

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	MXN Currency = "MXN"
)

func (c Currency) ToString() string {
	return string(c)
}

func NewCurrency(currency string) (Currency, error) {
	switch Currency(currency) {
	case USD:
		return USD, nil
	case EUR:
		return EUR, nil
	case MXN:
		return MXN, nil
	default:
		return "", ErrInvalidCurrency
	}
}

type CurrencyRateStatus string

const (
	PENDING CurrencyRateStatus = "PENDING"
	SUCCESS CurrencyRateStatus = "SUCCESS"
	ERROR   CurrencyRateStatus = "ERROR"
)

type CurrencyRate struct {
	Id       uuid.UUID
	RateTime *time.Time
	Status   CurrencyRateStatus
	Base     Currency
	Target   Currency
	Rate     *float32
}

type UpdateCurrencyRateMsg struct {
	RequestId uuid.UUID
	Base      Currency
	Target    Currency
}
