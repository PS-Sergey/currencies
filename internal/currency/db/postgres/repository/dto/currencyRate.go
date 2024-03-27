package dto

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

type CurrencyRateStatus string

const (
	PENDING CurrencyRateStatus = "PENDING"
	SUCCESS CurrencyRateStatus = "SUCCESS"
	ERROR   CurrencyRateStatus = "ERROR"
)

type CurrencyRate struct {
	Id       uuid.UUID          `db:"id"`
	RateTime *time.Time         `db:"rate_time"`
	Status   CurrencyRateStatus `db:"status"`
	Base     Currency           `db:"base"`
	Target   Currency           `db:"target"`
	Rate     *float32           `db:"rate"`
}
