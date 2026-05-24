package main

import (
	"github.com/d1slike/omlx-exporter/config"
)

type Config struct {
	Address  string `validate:"required"`
	Omlx     config.OmlxExporter
	LogLevel string
}
