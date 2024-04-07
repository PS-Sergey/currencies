package tracer

type Config struct {
	Host        string  `yaml:"host"`
	Ratio       float64 `yaml:"ratio"`
	ServiceName string  `yaml:"serviceName"`
}
