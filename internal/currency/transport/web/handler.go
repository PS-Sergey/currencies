package web

import (
	"currencies/api"
	"currencies/internal/currency"
	"currencies/internal/currency/tech"
	"currencies/internal/currency/types"
	"currencies/internal/currency/types/mappers"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type RateHandler struct {
	currencyRateService currency.RateService
	logger              *zap.Logger
}

func NewCurrencyRateHandler(currencyRateService currency.RateService, logger *zap.Logger) *RateHandler {
	return &RateHandler{
		currencyRateService: currencyRateService,
		logger:              logger,
	}
}

func (h *RateHandler) NewCurrencyRate(w http.ResponseWriter, r *http.Request, params api.NewCurrencyRateParams) {
	baseCurrency, targetCurrency, err := getCurrenciesFromPair(params.Pair)
	if err != nil {
		h.logErrAndSend(w, err)
		return
	}

	currencyRateId, err := h.currencyRateService.NewCurrencyRate(r.Context(), baseCurrency, targetCurrency)
	if err != nil {
		h.logErrAndSend(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(api.CurrencyRateIdResponse{CurrencyRateId: currencyRateId})
	if err != nil {
		h.logErrAndSend(w, err)
	}
}

func (h *RateHandler) GetLastCurrencyRate(w http.ResponseWriter, r *http.Request, base api.Currency, target api.Currency) {
	ctx, span := tech.StartSpan(r.Context(), "get last currency rate")
	defer span.End()

	baseCurrency, err := mappers.ApiCurrencyToDomain(base)
	if err != nil {
		h.logErrAndSend(w, err)
		return
	}

	targetCurrency, err := mappers.ApiCurrencyToDomain(target)
	if err != nil {
		h.logErrAndSend(w, err)
		return
	}

	currencyRate, err := h.currencyRateService.GetLastCurrencyRate(ctx, baseCurrency, targetCurrency)
	if err != nil {
		h.logErrAndSend(w, err)
		return
	}

	currencyRateResponse, err := mappers.DomainCurrencyRateToAPI(currencyRate)
	if err != nil {
		h.logErrAndSend(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(currencyRateResponse)
	if err != nil {
		h.logErrAndSend(w, err)
	}
}

func (h *RateHandler) GetCurrencyRate(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	currencyRate, err := h.currencyRateService.GetCurrencyRateById(r.Context(), id)
	if err != nil {
		if errors.Is(err, types.ErrCurrencyRateNotFound) {
			h.logErrAndSend(w, err)
			return
		}
		h.logErrAndSend(w, err)
		return
	}

	currencyRateResponse, err := mappers.DomainCurrencyRateToAPI(currencyRate)
	if err != nil {
		h.logErrAndSend(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(currencyRateResponse)
	if err != nil {
		h.logErrAndSend(w, err)
	}
}

func (h *RateHandler) logErrAndSend(w http.ResponseWriter, err error) {
	h.logger.Error("got error", zap.Error(err))
	var responseErr api.Error

	switch {
	case errors.Is(err, types.ErrInvalidCurrencyPair):
		responseErr = api.Error{Code: http.StatusBadRequest, Message: types.ErrInvalidCurrencyPair.Error()}
	case errors.Is(err, types.ErrInvalidCurrency):
		responseErr = api.Error{Code: http.StatusBadRequest, Message: types.ErrInvalidCurrency.Error()}
	case errors.Is(err, types.ErrCurrencyRateNotFound):
		responseErr = api.Error{Code: http.StatusNotFound, Message: types.ErrCurrencyRateNotFound.Error()}
	default:
		responseErr = api.Error{Code: http.StatusInternalServerError, Message: "server error"}
	}

	w.WriteHeader(int(responseErr.Code))
	err = json.NewEncoder(w).Encode(responseErr)
	if err != nil {
		h.logger.Error("got error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getCurrenciesFromPair(currencyPair string) (baseCurrency types.Currency, targetCurrency types.Currency, err error) {
	const currencyPairSep = "/"

	base, target, found := strings.Cut(currencyPair, currencyPairSep)
	if !found {
		err = types.ErrInvalidCurrencyPair
		return
	}

	baseCurrency, err = types.NewCurrency(base)
	if err != nil {
		return
	}

	targetCurrency, err = types.NewCurrency(target)
	if err != nil {
		return
	}

	return
}
