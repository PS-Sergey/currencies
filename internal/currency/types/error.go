package types

import "errors"

var ErrInvalidCurrencyPair = errors.New("currency pair invalid format")
var ErrInvalidCurrency = errors.New("currency invalid format")
var ErrCurrencyRateNotFound = errors.New("currency rate not found")
var ErrServerError = errors.New("application error")
