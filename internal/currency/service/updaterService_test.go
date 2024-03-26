package service

import (
	"context"
	"currencies/internal/currency"
	mock_currency "currencies/internal/currency/mocks"
	"currencies/internal/currency/types"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCurrencyRateUpdateService_updateCurrencyRateTx(t *testing.T) {
	id := uuid.New()
	base := types.EUR
	target := types.USD
	rate := float32(1.1)
	currencyRate := types.CurrencyRate{
		Id:     id,
		Status: types.PENDING,
		Base:   base,
		Target: target,
	}
	successResultCurrencyRate := types.CurrencyRate{
		Id:     id,
		Status: types.SUCCESS,
		Base:   base,
		Target: target,
		Rate:   &rate,
	}
	errorResultCurrencyRate := types.CurrencyRate{
		Id:     id,
		Status: types.ERROR,
		Base:   base,
		Target: target,
	}
	updateCurrencyRateMsg := types.UpdateCurrencyRateMsg{
		RequestId: id,
		Base:      base,
		Target:    target,
	}

	type args struct {
		makeTransaction func(ctrl *gomock.Controller) currency.Transaction
		msg             types.UpdateCurrencyRateMsg
	}
	tests := []struct {
		name                   string
		makeExchangeRateClient func(ctrl *gomock.Controller) currency.ExchangeRateClient
		args                   args
		wantErr                bool
	}{
		{
			name: "Success update currency rate",
			makeExchangeRateClient: func(ctrl *gomock.Controller) currency.ExchangeRateClient {
				mock := mock_currency.NewMockExchangeRateClient(ctrl)
				mock.EXPECT().GetCurrencyRate(gomock.Any(), base, target).
					Return(rate, nil).Times(1)
				return mock
			},
			args: args{
				makeTransaction: func(ctrl *gomock.Controller) currency.Transaction {
					repoMock := mock_currency.NewMockRateRepository(ctrl)
					repoMock.EXPECT().GetCurrencyRateById(gomock.Any(), id).
						Return(currencyRate, nil).Times(1)
					repoMock.EXPECT().UpdateCurrencyRate(gomock.Any(), gomock.All(
						matchCurrencyRates(successResultCurrencyRate),
					)).
						Return(nil).Times(1)

					mock := mock_currency.NewMockTransaction(ctrl)
					mock.EXPECT().TxRateRepository().Return(repoMock)
					return mock
				},
				msg: updateCurrencyRateMsg,
			},
			wantErr: false,
		},
		{
			name: "Error get currency rate from repository",
			makeExchangeRateClient: func(ctrl *gomock.Controller) currency.ExchangeRateClient {
				mock := mock_currency.NewMockExchangeRateClient(ctrl)
				mock.EXPECT().GetCurrencyRate(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(rate, nil).Times(0)
				return mock
			},
			args: args{
				makeTransaction: func(ctrl *gomock.Controller) currency.Transaction {
					repoMock := mock_currency.NewMockRateRepository(ctrl)
					repoMock.EXPECT().GetCurrencyRateById(gomock.Any(), id).
						Return(types.CurrencyRate{}, errSomething).Times(1)
					repoMock.EXPECT().UpdateCurrencyRate(gomock.Any(), gomock.Any()).
						Return(nil).Times(0)

					mock := mock_currency.NewMockTransaction(ctrl)
					mock.EXPECT().TxRateRepository().Return(repoMock)
					return mock
				},
				msg: updateCurrencyRateMsg,
			},
			wantErr: true,
		},
		{
			name: "Error get currency rate from client",
			makeExchangeRateClient: func(ctrl *gomock.Controller) currency.ExchangeRateClient {
				mock := mock_currency.NewMockExchangeRateClient(ctrl)
				mock.EXPECT().GetCurrencyRate(gomock.Any(), base, target).
					Return(float32(0), errSomething).Times(1)
				return mock
			},
			args: args{
				makeTransaction: func(ctrl *gomock.Controller) currency.Transaction {
					repoMock := mock_currency.NewMockRateRepository(ctrl)
					repoMock.EXPECT().GetCurrencyRateById(gomock.Any(), id).
						Return(currencyRate, nil).Times(1)
					repoMock.EXPECT().UpdateCurrencyRate(gomock.Any(), gomock.All(
						matchCurrencyRates(errorResultCurrencyRate),
					)).
						Return(nil).Times(1)

					mock := mock_currency.NewMockTransaction(ctrl)
					mock.EXPECT().TxRateRepository().Return(repoMock)
					return mock
				},
				msg: updateCurrencyRateMsg,
			},
			wantErr: false,
		},
		{
			name: "Error update currency rate in repository",
			makeExchangeRateClient: func(ctrl *gomock.Controller) currency.ExchangeRateClient {
				mock := mock_currency.NewMockExchangeRateClient(ctrl)
				mock.EXPECT().GetCurrencyRate(gomock.Any(), base, target).
					Return(rate, nil).Times(1)
				return mock
			},
			args: args{
				makeTransaction: func(ctrl *gomock.Controller) currency.Transaction {
					repoMock := mock_currency.NewMockRateRepository(ctrl)
					repoMock.EXPECT().GetCurrencyRateById(gomock.Any(), id).
						Return(currencyRate, nil).Times(1)
					repoMock.EXPECT().UpdateCurrencyRate(gomock.Any(), gomock.Any()).
						Return(errSomething).Times(1)

					mock := mock_currency.NewMockTransaction(ctrl)
					mock.EXPECT().TxRateRepository().Return(repoMock)
					return mock
				},
				msg: updateCurrencyRateMsg,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := CurrencyRateUpdateService{
				exchangeRateClient: tt.makeExchangeRateClient(ctrl),
				logger:             zap.NewExample(),
			}

			err := s.updateCurrencyRateTx(context.Background(), tt.args.makeTransaction(ctrl), tt.args.msg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func matchCurrencyRates(currencyRate types.CurrencyRate) gomock.Matcher {
	return &currencyRateMatcher{
		baseCurrencyRate: currencyRate,
	}
}

type currencyRateMatcher struct {
	baseCurrencyRate types.CurrencyRate
}

func (c *currencyRateMatcher) Matches(x interface{}) bool {
	currencyRate, ok := x.(types.CurrencyRate)
	if !ok {
		return false
	}

	return c.baseCurrencyRate.Id == currencyRate.Id &&
		c.baseCurrencyRate.Base == currencyRate.Base &&
		c.baseCurrencyRate.Target == currencyRate.Target &&
		c.baseCurrencyRate.Status == currencyRate.Status &&
		(c.baseCurrencyRate.Rate != nil && currencyRate.Rate != nil && *c.baseCurrencyRate.Rate == *currencyRate.Rate ||
			c.baseCurrencyRate.Rate == nil && currencyRate.Rate == nil)
}

func (c *currencyRateMatcher) String() string {
	return fmt.Sprintf("is matched to %v", c.baseCurrencyRate)
}
