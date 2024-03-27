package service

import (
	"context"
	"currencies/internal/currency"
	"currencies/internal/currency/db/postgres"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type TransactionManager struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewTransactionManager(db *sqlx.DB, logger *zap.Logger) *TransactionManager {
	return &TransactionManager{
		db:     db,
		logger: logger,
	}
}

func (tm *TransactionManager) BeginTx(ctx context.Context, txOptions *sql.TxOptions, f func(ctx context.Context, t currency.Transaction) error) error {
	tx, err := tm.db.BeginTxx(ctx, txOptions)
	if err != nil {
		return errors.Wrap(err, "start transaction")
	}

	transaction := postgres.NewTransaction(tx)
	err = f(ctx, transaction)

	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			tm.logger.Error("failed to rollback tx")
		}
		return errors.Wrap(err, "run function in tx")
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit tx")
	}

	return nil
}
