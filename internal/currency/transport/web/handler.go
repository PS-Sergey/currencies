package web

import (
	"currencies/api"
	"currencies/internal/currency"
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
		h.logErrAndSend(w, http.StatusBadRequest, err)
		return
	}

	currencyRateId, err := h.currencyRateService.NewCurrencyRate(r.Context(), baseCurrency, targetCurrency)
	if err != nil {
		h.logErrAndSend(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(api.CurrencyRateIdResponse{CurrencyRateId: currencyRateId})
}

func (h *RateHandler) GetLastCurrencyRate(w http.ResponseWriter, r *http.Request, base api.Currency, target api.Currency) {
	baseCurrency, err := mappers.ApiCurrencyToDomain(base)
	if err != nil {
		h.logErrAndSend(w, http.StatusBadRequest, err)
		return
	}

	targetCurrency, err := mappers.ApiCurrencyToDomain(target)
	if err != nil {
		h.logErrAndSend(w, http.StatusBadRequest, err)
		return
	}

	currencyRate, err := h.currencyRateService.GetLastCurrencyRate(r.Context(), baseCurrency, targetCurrency)
	if err != nil {
		if errors.Is(err, types.ErrCurrencyRateNotFound) {
			h.logErrAndSend(w, http.StatusNotFound, err)
			return
		}
		h.logErrAndSend(w, http.StatusInternalServerError, err)
		return
	}

	currencyRateResponse, err := mappers.DomainCurrencyRateToAPI(currencyRate)
	if err != nil {
		h.logErrAndSend(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currencyRateResponse)
}

func (h *RateHandler) GetCurrencyRate(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	currencyRate, err := h.currencyRateService.GetCurrencyRateById(r.Context(), id)
	if err != nil {
		if errors.Is(err, types.ErrCurrencyRateNotFound) {
			h.logErrAndSend(w, http.StatusNotFound, err)
			return
		}
		h.logErrAndSend(w, http.StatusInternalServerError, err)
		return
	}

	currencyRateResponse, err := mappers.DomainCurrencyRateToAPI(currencyRate)
	if err != nil {
		h.logErrAndSend(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currencyRateResponse)
}

func (h *RateHandler) logErrAndSend(w http.ResponseWriter, code int, err error) {
	h.logger.Error("got error", zap.Error(err))

	petErr := api.Error{
		Code:    int32(code),
		Message: err.Error(),
	}

	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(petErr)
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
