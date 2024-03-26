package repository

import (
	"context"
	"currencies/internal/currency/db"
	"currencies/internal/currency/db/postgres/repository/dto"
	"currencies/internal/currency/db/postgres/repository/dto/mappers"
	"currencies/internal/currency/types"
	"database/sql"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type RateRepository struct {
	db db.Executor
}

func NewRateRepository(db db.Executor) *RateRepository {
	return &RateRepository{
		db: db,
	}
}

func (r *RateRepository) SaveCurrencyRate(ctx context.Context, currencyRate types.CurrencyRate) error {
	dbCurrencyRate, err := mappers.DomainCurrencyRateToDB(currencyRate)
	if err != nil {
		return errors.Wrap(err, "mapping domain currency rate to DB")
	}

	_, err = r.db.NamedExecContext(ctx, insertCurrencyRateQuery, dbCurrencyRate)
	if err != nil {
		return errors.Wrap(err, "insert currency rate")
	}

	return nil
}

func (r *RateRepository) GetCurrencyRateById(ctx context.Context, id uuid.UUID) (types.CurrencyRate, error) {
	var currencyRate dto.CurrencyRate

	getCurrencyRateBindQuery, getCurrencyRateArgs, err := r.db.BindNamed(getCurrencyRateByIdQuery, map[string]any{"id": id})
	if err != nil {
		return types.CurrencyRate{}, errors.Wrap(err, "get currency rate by id binds query")
	}

	err = r.db.GetContext(ctx, &currencyRate, getCurrencyRateBindQuery, getCurrencyRateArgs...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.CurrencyRate{}, types.ErrCurrencyRateNotFound
		}
		return types.CurrencyRate{}, errors.Wrap(err, "query get currency rate by id")
	}

	return mappers.DBCurrencyRateToDomain(currencyRate)
}

func (r *RateRepository) GetLastCurrencyRate(ctx context.Context, current types.Currency, target types.Currency) (types.CurrencyRate, error) {
	var currencyRate dto.CurrencyRate

	getLastCurrencyRateBindQuery, getLastCurrencyRateArgs, err := r.db.BindNamed(getLastCurrencyRateQuery, map[string]any{
		"base":   current,
		"target": target,
	})
	if err != nil {
		return types.CurrencyRate{}, errors.Wrap(err, "get last currency rate binds query")
	}

	err = r.db.GetContext(ctx, &currencyRate, getLastCurrencyRateBindQuery, getLastCurrencyRateArgs...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.CurrencyRate{}, types.ErrCurrencyRateNotFound
		}
		return types.CurrencyRate{}, errors.Wrap(err, "query get last currency rate")
	}

	return mappers.DBCurrencyRateToDomain(currencyRate)
}

func (r *RateRepository) UpdateCurrencyRate(ctx context.Context, rate types.CurrencyRate) error {
	dbCurrencyRate, err := mappers.DomainCurrencyRateToDB(rate)
	if err != nil {
		return errors.Wrap(err, "mapping domain currency rate to DB")
	}

	_, err = r.db.NamedExecContext(ctx, updateCurrencyRateQuery, dbCurrencyRate)
	if err != nil {
		return errors.Wrap(err, "update currency rate")
	}

	return nil
}

func (r *RateRepository) WithTx(ctx context.Context) {

}
