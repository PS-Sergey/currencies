package service

import (
	"context"
	"currencies/internal/currency"
	mock_currency "currencies/internal/currency/mocks"
	"currencies/internal/currency/types"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var errSomething = errors.New("something went wrong")

func TestCurrencyRateService_NewCurrencyRate(t *testing.T) {
	id := uuid.New()
	base := types.EUR
	target := types.USD
	newCurrencyRate := types.CurrencyRate{
		Id:     id,
		Status: types.PENDING,
		Base:   base,
		Target: target,
	}
	updateCurrencyRateMsg := types.UpdateCurrencyRateMsg{
		RequestId: id,
		Base:      base,
		Target:    target,
	}

	type args struct {
		base   types.Currency
		target types.Currency
	}
	tests := []struct {
		name                   string
		makeRepository         func(ctrl *gomock.Controller) currency.RateRepository
		updateCurrencyRateChan chan types.UpdateCurrencyRateMsg
		makeUuidGenerator      func(ctrl *gomock.Controller) currency.UUIDGenerator
		args                   args
		wantId                 uuid.UUID
		wantUpdateMsg          types.UpdateCurrencyRateMsg
		wantErr                bool
	}{
		{
			name: "Success save new currency rate",
			makeRepository: func(ctrl *gomock.Controller) currency.RateRepository {
				mock := mock_currency.NewMockRateRepository(ctrl)
				mock.EXPECT().SaveCurrencyRate(gomock.Any(), newCurrencyRate).
					Return(nil).Times(1)
				return mock
			},
			updateCurrencyRateChan: make(chan types.UpdateCurrencyRateMsg, 1),
			makeUuidGenerator: func(ctrl *gomock.Controller) currency.UUIDGenerator {
				mock := mock_currency.NewMockUUIDGenerator(ctrl)
				mock.EXPECT().Generate().Return(id, nil).Times(1)
				return mock
			},
			args: args{
				base:   base,
				target: target,
			},
			wantId:        id,
			wantUpdateMsg: updateCurrencyRateMsg,
			wantErr:       false,
		},
		{
			name: "Error generate id",
			makeRepository: func(ctrl *gomock.Controller) currency.RateRepository {
				mock := mock_currency.NewMockRateRepository(ctrl)
				mock.EXPECT().SaveCurrencyRate(gomock.Any(), gomock.Any()).Times(0)
				return mock
			},
			updateCurrencyRateChan: make(chan types.UpdateCurrencyRateMsg, 1),
			makeUuidGenerator: func(ctrl *gomock.Controller) currency.UUIDGenerator {
				mock := mock_currency.NewMockUUIDGenerator(ctrl)
				mock.EXPECT().Generate().Return(uuid.UUID{}, errSomething).Times(1)
				return mock
			},
			args: args{
				base:   base,
				target: target,
			},
			wantId:  uuid.UUID{},
			wantErr: true,
		},
		{
			name: "Error save nuw currency rate in repository",
			makeRepository: func(ctrl *gomock.Controller) currency.RateRepository {
				mock := mock_currency.NewMockRateRepository(ctrl)
				mock.EXPECT().SaveCurrencyRate(gomock.Any(), newCurrencyRate).
					Return(errSomething).Times(1)
				return mock
			},
			updateCurrencyRateChan: make(chan types.UpdateCurrencyRateMsg, 1),
			makeUuidGenerator: func(ctrl *gomock.Controller) currency.UUIDGenerator {
				mock := mock_currency.NewMockUUIDGenerator(ctrl)
				mock.EXPECT().Generate().Return(id, nil).Times(1)
				return mock
			},
			args: args{
				base:   base,
				target: target,
			},
			wantId:  uuid.UUID{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := CurrencyRateService{
				currencyRateRepository: tt.makeRepository(ctrl),
				updateCurrencyRateChan: tt.updateCurrencyRateChan,
				uuidGenerator:          tt.makeUuidGenerator(ctrl),
			}

			got, err := s.NewCurrencyRate(context.Background(), tt.args.base, tt.args.target)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				gotMsg := <-tt.updateCurrencyRateChan
				require.NoError(t, err)
				require.Equal(t, tt.wantId, got)
				require.Equal(t, tt.wantUpdateMsg, gotMsg)
			}
		})
	}
}

func TestCurrencyRateService_GetCurrencyRateById(t *testing.T) {
	id := uuid.New()
	currencyRate := types.CurrencyRate{
		Id:     id,
		Status: types.PENDING,
		Base:   types.EUR,
		Target: types.USD,
	}

	tests := []struct {
		name           string
		makeRepository func(ctrl *gomock.Controller) currency.RateRepository
		currencyRateId uuid.UUID
		want           types.CurrencyRate
		wantErr        bool
	}{
		{
			name: "Success get currency rate by id",
			makeRepository: func(ctrl *gomock.Controller) currency.RateRepository {
				mock := mock_currency.NewMockRateRepository(ctrl)
				mock.EXPECT().GetCurrencyRateById(gomock.Any(), id).
					Return(currencyRate, nil).Times(1)
				return mock
			},
			currencyRateId: id,
			want:           currencyRate,
			wantErr:        false,
		},
		{
			name: "Error get currency rate by id in repo",
			makeRepository: func(ctrl *gomock.Controller) currency.RateRepository {
				mock := mock_currency.NewMockRateRepository(ctrl)
				mock.EXPECT().GetCurrencyRateById(gomock.Any(), id).
					Return(types.CurrencyRate{}, errSomething).Times(1)
				return mock
			},
			currencyRateId: id,
			want:           types.CurrencyRate{},
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := CurrencyRateService{
				currencyRateRepository: tt.makeRepository(ctrl),
			}

			got, err := s.GetCurrencyRateById(context.Background(), tt.currencyRateId)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestCurrencyRateService_GetLastCurrencyRate(t *testing.T) {
	base := types.EUR
	target := types.USD
	currencyRate := types.CurrencyRate{
		Id:     uuid.New(),
		Status: types.ERROR,
		Base:   base,
		Target: target,
	}

	type args struct {
		base   types.Currency
		target types.Currency
	}
	tests := []struct {
		name           string
		makeRepository func(ctrl *gomock.Controller) currency.RateRepository
		args           args
		want           types.CurrencyRate
		wantErr        bool
	}{
		{
			name: "Success get last currency rate",
			makeRepository: func(ctrl *gomock.Controller) currency.RateRepository {
				mock := mock_currency.NewMockRateRepository(ctrl)
				mock.EXPECT().GetLastCurrencyRate(gomock.Any(), base, target).
					Return(currencyRate, nil).Times(1)
				return mock
			},
			args: args{
				base:   base,
				target: target,
			},
			want:    currencyRate,
			wantErr: false,
		},
		{
			name: "Error get last currency rate in repository",
			makeRepository: func(ctrl *gomock.Controller) currency.RateRepository {
				mock := mock_currency.NewMockRateRepository(ctrl)
				mock.EXPECT().GetLastCurrencyRate(gomock.Any(), base, target).
					Return(types.CurrencyRate{}, errSomething).Times(1)
				return mock
			},
			args: args{
				base:   base,
				target: target,
			},
			want:    types.CurrencyRate{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := CurrencyRateService{
				currencyRateRepository: tt.makeRepository(ctrl),
			}

			got, err := s.GetLastCurrencyRate(context.Background(), tt.args.base, tt.args.target)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
