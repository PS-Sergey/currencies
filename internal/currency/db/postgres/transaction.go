package postgres

import (
	"currencies/internal/currency"
	"currencies/internal/currency/db/postgres/repository"

	"github.com/jmoiron/sqlx"
)

type Transaction struct {
	tx *sqlx.Tx
}

func NewTransaction(tx *sqlx.Tx) *Transaction {
	return &Transaction{
		tx: tx,
	}
}

func (t *Transaction) TxRateRepository() currency.RateRepository {
	return repository.NewRateRepository(t.tx)
}
