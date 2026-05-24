package main

import (
	"context"
	stdlog "log"
	"time"

	"github.com/pkg/errors"
	"github.com/txix-open/isp-kit/app"
	ispconfig "github.com/txix-open/isp-kit/config"
	"github.com/txix-open/isp-kit/infra"
	"github.com/txix-open/isp-kit/log"
	"github.com/txix-open/isp-kit/metrics"
	"github.com/txix-open/isp-kit/shutdown"
	"github.com/txix-open/isp-kit/validator"
	"go.uber.org/zap/zapcore"

	"github.com/d1slike/omlx-exporter/client"
	"github.com/d1slike/omlx-exporter/config"
	"github.com/d1slike/omlx-exporter/exporter"
)

func main() {
	application, err := app.New(
		app.WithConfigOptions(
			ispconfig.WithExtraSource(ispconfig.NewYamlConfig("config.yaml")),
			ispconfig.WithValidator(validator.Default),
		),
	)
	if err != nil {
		stdlog.Fatal(err)
	}

	ctx := application.Context()
	logger := application.Logger()
	cfg := config.Config{}
	err = application.Config().Read(&cfg)
	if err != nil {
		logger.Fatal(ctx, errors.WithMessage(err, "read config"))
	}
	logLevel, err := zapcore.ParseLevel(cfg.LogLevel)
	if err != nil {
		logger.Fatal(ctx, "failed to parse log level")
	}
	logger.SetLevel(logLevel)

	srv := infra.NewServer()
	srv.Handle("GET /metrics", metrics.DefaultRegistry.MetricsHandler())
	srv.Handle("GET /metrics/description", metrics.DefaultRegistry.MetricsDescriptionHandler())

	interval := time.Duration(cfg.Omlx.ScrapeIntervalInSec) * time.Second
	omlxClient := client.New(cfg.Omlx.TargetUrl, cfg.Omlx.ApiKey, logger)
	err = omlxClient.RefreshSession(ctx)
	if err != nil {
		logger.Fatal(ctx, errors.WithMessage(err, "refresh omlx admin session"))
	}
	scraper := exporter.New(omlxClient, interval, logger)

	application.AddRunners(app.RunnerFunc(func(ctx context.Context) error {
		logger.Info(ctx, "omlx exporter starting",
			log.String("target", cfg.Omlx.TargetUrl),
			log.String("interval", interval.String()))
		return scraper.Run(ctx)
	}))
	application.AddClosers(app.CloserFunc(func() error {
		logger.Info(ctx, "omlx exporter stopping")
		scraper.Stop()
		return nil
	}))

	application.AddRunners(app.RunnerFunc(func(ctx context.Context) error {
		logger.Info(ctx, "listen and serve", log.String("address", cfg.Address))
		err := srv.ListenAndServe(cfg.Address)
		return errors.WithMessage(err, "listen and serve")
	}))
	application.AddClosers(app.CloserFunc(func() error {
		return srv.Shutdown()
	}))

	shutdown.On(func() {
		logger.Info(ctx, "shutdown start")
		application.Close()
		logger.Info(ctx, "shutdown done")
	})

	err = application.Run()
	if err != nil {
		logger.Fatal(ctx, errors.WithMessage(err, "run application"))
	}
}
