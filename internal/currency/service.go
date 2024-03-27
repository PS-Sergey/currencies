package currency

import (
	"context"
	"currencies/internal/currency/types"
	"database/sql"

	"github.com/google/uuid"
)

//go:generate mockgen -source service.go -destination mocks/service.go
type RateService interface {
	NewCurrencyRate(ctx context.Context, base types.Currency, target types.Currency) (uuid.UUID, error)
	GetCurrencyRateById(ctx context.Context, currencyRateId uuid.UUID) (types.CurrencyRate, error)
	GetLastCurrencyRate(ctx context.Context, base types.Currency, target types.Currency) (types.CurrencyRate, error)
}

type RateUpdaterService interface {
	UpdateCurrencyRate(ctx context.Context, msg types.UpdateCurrencyRateMsg)
}

type UUIDGenerator interface {
	Generate() (uuid.UUID, error)
}

type TransactionManager interface {
	BeginTx(ctx context.Context, txOptions *sql.TxOptions, f func(ctx context.Context, t Transaction) error) error
}
