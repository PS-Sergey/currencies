package currency

import (
	"context"
	"currencies/internal/currency/types"

	"github.com/google/uuid"
)

//go:generate mockgen -source repository.go -destination mocks/repository.go
type RateRepository interface {
	SaveCurrencyRate(ctx context.Context, currencyRate types.CurrencyRate) error
	GetCurrencyRateById(ctx context.Context, id uuid.UUID) (types.CurrencyRate, error)
	GetLastCurrencyRate(ctx context.Context, base types.Currency, target types.Currency) (types.CurrencyRate, error)
	UpdateCurrencyRate(ctx context.Context, rate types.CurrencyRate) error
}
