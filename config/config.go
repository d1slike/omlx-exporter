package config

type OmlxExporter struct {
	TargetUrl           string `validate:"required"`
	ApiKey              string `validate:"required"`
	ScrapeIntervalInSec int    `validate:"required"`
}
