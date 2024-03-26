// Code generated by MockGen. DO NOT EDIT.
// Source: clients.go

// Package mock_currency is a generated GoMock package.
package mock_currency

import (
	context "context"
	types "currencies/internal/currency/types"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockExchangeRateClient is a mock of ExchangeRateClient interface.
type MockExchangeRateClient struct {
	ctrl     *gomock.Controller
	recorder *MockExchangeRateClientMockRecorder
}

// MockExchangeRateClientMockRecorder is the mock recorder for MockExchangeRateClient.
type MockExchangeRateClientMockRecorder struct {
	mock *MockExchangeRateClient
}

// NewMockExchangeRateClient creates a new mock instance.
func NewMockExchangeRateClient(ctrl *gomock.Controller) *MockExchangeRateClient {
	mock := &MockExchangeRateClient{ctrl: ctrl}
	mock.recorder = &MockExchangeRateClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExchangeRateClient) EXPECT() *MockExchangeRateClientMockRecorder {
	return m.recorder
}

// GetCurrencyRate mocks base method.
func (m *MockExchangeRateClient) GetCurrencyRate(ctx context.Context, current, target types.Currency) (float32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrencyRate", ctx, current, target)
	ret0, _ := ret[0].(float32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrencyRate indicates an expected call of GetCurrencyRate.
func (mr *MockExchangeRateClientMockRecorder) GetCurrencyRate(ctx, current, target interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrencyRate", reflect.TypeOf((*MockExchangeRateClient)(nil).GetCurrencyRate), ctx, current, target)
}
