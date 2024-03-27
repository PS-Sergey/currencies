package currency

import (
	"context"
	"currencies/internal/currency/types"
)

//go:generate mockgen -source clients.go -destination mocks/clients.go
type ExchangeRateClient interface {
	GetCurrencyRate(ctx context.Context, base types.Currency, target types.Currency) (float32, error)
}
