package main

import (
	"context"
	log2 "log"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/app"
	"github.com/txix-open/isp-kit/config"
	"github.com/txix-open/isp-kit/infra"
	"github.com/txix-open/isp-kit/log"
	"github.com/txix-open/isp-kit/metrics"
	"github.com/txix-open/isp-kit/shutdown"
	"github.com/txix-open/isp-kit/validator"
)

type Config struct {
	Address string `validate:"required"`
}

func main() {
	applicaion, err := app.New(
		app.WithConfigOptions(
			config.WithExtraSource(config.NewYamlConfig("config.yaml")),
			config.WithValidator(validator.Default),
		),
	)
	if err != nil {
		log2.Fatal(err)
	}

	ctx := applicaion.Context()
	logger := applicaion.Logger()
	cfg := Config{}
	err = applicaion.Config().Read(&cfg)
	if err != nil {
		logger.Fatal(ctx, errors.WithMessage(err, "read config"))
	}

	srv := infra.NewServer()
	srv.Handle("GET /metrics", metrics.DefaultRegistry.MetricsHandler())
	applicaion.AddRunners(app.RunnerFunc(func(ctx context.Context) error {
		logger.Info(ctx, "listen and serve", log.String("address", cfg.Address))
		err := srv.ListenAndServe(cfg.Address)
		return errors.WithMessage(err, "listen and serve")
	}))
	applicaion.AddClosers(app.CloserFunc(func() error {
		return srv.Shutdown()
	}))

	shutdown.On(func() {
		logger.Info(ctx, "shutdown start")
		applicaion.Close()
		logger.Info(ctx, "shutdown done")
	})

	err = applicaion.Run()
	if err != nil {
		logger.Fatal(ctx, errors.WithMessage(err, "run application"))
	}
}
