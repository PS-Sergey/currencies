package config

import (
	"currencies/cmd/application"
	"currencies/pkg/clients/exchangeRate"
	"currencies/pkg/db/postgres"
)

type Config struct {
	Server             *application.Config  `yaml:"server"`
	DB                 *postgres.Config     `yaml:"db"`
	ExchangeRateClient *exchangeRate.Config `yaml:"exchangeRateClient"`
}
