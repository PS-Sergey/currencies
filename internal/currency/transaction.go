package currency

//go:generate mockgen -source transaction.go -destination mocks/transaction.go
type Transaction interface {
	TxRateRepository() RateRepository
}
