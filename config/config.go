package config

type Config struct {
	Address  string `validate:"required"`
	Omlx     OmlxExporter
	LogLevel string
}

type OmlxExporter struct {
	TargetUrl           string `validate:"required"`
	ApiKey              string `validate:"required"`
	ScrapeIntervalInSec int    `validate:"required"`
}
