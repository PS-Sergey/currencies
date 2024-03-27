package exchangeRate

import "time"

type Config struct {
	BaseUrl string        `yaml:"baseUrl"`
	Timeout time.Duration `yaml:"timeout"`
	ApiKey  string        `yaml:"apiKey"`
}
