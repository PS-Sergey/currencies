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
	h.logger.Info("NewCurrencyRate request with params:", zap.String("pair", params.Pair))
	baseCurrency, targetCurrency, err := getCurrenciesFromPair(params.Pair)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	currencyRateId, err := h.currencyRateService.NewCurrencyRate(r.Context(), baseCurrency, targetCurrency)
	if err != nil {
		sendError(w, http.StatusInternalServerError, types.ErrServerError.Error())
		return
	}

	h.logger.Info("NewCurrencyRate response:", zap.String("new currency rate id", currencyRateId.String()))
	setContentTypeAndStatus(w, http.StatusCreated)
	json.NewEncoder(w).Encode(api.CurrencyRateIdResponse{CurrencyRateId: currencyRateId})
}

func (h *RateHandler) GetLastCurrencyRate(w http.ResponseWriter, r *http.Request, base api.Currency, target api.Currency) {
	baseCurrency, err := mappers.ApiCurrencyToDomain(base)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid format for base currency")
		return
	}

	targetCurrency, err := mappers.ApiCurrencyToDomain(target)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid format for target currency")
		return
	}

	currencyRate, err := h.currencyRateService.GetLastCurrencyRate(r.Context(), baseCurrency, targetCurrency)
	if err != nil {
		if errors.Is(err, types.ErrCurrencyRateNotFound) {
			sendError(w, http.StatusNotFound, types.ErrCurrencyRateNotFound.Error())
			return
		}
		sendError(w, http.StatusInternalServerError, types.ErrServerError.Error())
		return
	}

	currencyRateResponse, err := mappers.DomainCurrencyRateToAPI(currencyRate)
	if err != nil {
		sendError(w, http.StatusInternalServerError, types.ErrServerError.Error())
		return
	}

	setContentTypeAndStatus(w, http.StatusOK)
	json.NewEncoder(w).Encode(currencyRateResponse)
}

func (h *RateHandler) GetCurrencyRate(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	currencyRate, err := h.currencyRateService.GetCurrencyRateById(r.Context(), id)
	if err != nil {
		if errors.Is(err, types.ErrCurrencyRateNotFound) {
			sendError(w, http.StatusNotFound, types.ErrCurrencyRateNotFound.Error())
			return
		}
		sendError(w, http.StatusInternalServerError, types.ErrServerError.Error())
		return
	}

	currencyRateResponse, err := mappers.DomainCurrencyRateToAPI(currencyRate)
	if err != nil {
		sendError(w, http.StatusInternalServerError, types.ErrServerError.Error())
		return
	}

	setContentTypeAndStatus(w, http.StatusOK)
	json.NewEncoder(w).Encode(currencyRateResponse)
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

func setContentTypeAndStatus(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
}

func sendError(w http.ResponseWriter, code int, message string) {
	petErr := api.Error{
		Code:    int32(code),
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(petErr)
}
