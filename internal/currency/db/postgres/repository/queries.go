package repository

import _ "embed"

var (
	//go:embed queries/insert_currency_rate.sql
	insertCurrencyRateQuery string

	//go:embed queries/get_currency_rate_by_id.sql
	getCurrencyRateByIdQuery string

	//go:embed queries/get_last_currency_rate.sql
	getLastCurrencyRateQuery string

	//go:embed queries/update_currency_rate.sql
	updateCurrencyRateQuery string
)
