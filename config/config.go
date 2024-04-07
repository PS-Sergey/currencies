package config

import (
	"currencies/cmd/application"
	"currencies/pkg/clients/exchangeRate"
	"currencies/pkg/db/postgres"
	"currencies/pkg/tracer"
)

type Config struct {
	Server             *application.Config  `yaml:"server"`
	DB                 *postgres.Config     `yaml:"db"`
	Tracer             *tracer.Config       `yaml:"tracer"`
	ExchangeRateClient *exchangeRate.Config `yaml:"exchangeRateClient"`
}
