package exchangeRate

import (
	"context"
	"currencies/internal/currency/types"

	"github.com/pkg/errors"
)

type ExchangeRequestClient interface {
	GetCurrencyRate(ctx context.Context, from string, to string) (float32, error)
}

type Adapter struct {
	c ExchangeRequestClient
}

func NewExchangeRateAdapter(c ExchangeRequestClient) *Adapter {
	return &Adapter{
		c: c,
	}
}

func (a *Adapter) GetCurrencyRate(ctx context.Context, base types.Currency, target types.Currency) (float32, error) {
	rate, err := a.c.GetCurrencyRate(ctx, base.ToString(), target.ToString())
	if err != nil {
		return 0, errors.Wrap(err, "get currency rate from client")
	}

	return rate, nil
}
