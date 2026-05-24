package exporter

import (
	"context"
	"time"

	"github.com/d1slike/omlx-exporter/client"
	"github.com/d1slike/omlx-exporter/metrics"
	"github.com/txix-open/isp-kit/log"
)

type Scraper struct {
	client   *client.Client
	interval time.Duration
	stopCh   chan struct{}
	logger   log.Logger
}

func New(c *client.Client, interval time.Duration, logger log.Logger) *Scraper {
	return &Scraper{
		client:   c,
		interval: interval,
		stopCh:   make(chan struct{}),
		logger:   logger,
	}
}

func (s *Scraper) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-s.stopCh:
			return nil
		case <-ticker.C:
			s.scrape(ctx)
		}
	}
}

func (s *Scraper) Stop() {
	close(s.stopCh)
}

func (s *Scraper) scrape(ctx context.Context) {
	start := time.Now()
	resp, err := s.client.FetchStats(ctx)
	if err != nil {
		metrics.RecordFailure()
		if s.client.NeedsSessionRefresh(err) {
			s.logger.Info(ctx, "refreshing session")
			if refreshErr := s.client.RefreshSession(ctx); refreshErr != nil {
				s.logger.Error(ctx, "failed to refresh session", log.Any("error", refreshErr))
				return
			}
			resp, err = s.client.FetchStats(ctx)
			if err != nil {
				s.logger.Error(ctx, "failed to fetch stats after session refresh", log.Any("error", err))
				return
			}
		}
		s.logger.Error(ctx, "failed to fetch stats", log.Any("error", err))
		return
	}
	s.updateStats(resp)
	s.updateActiveModels(resp)
	metrics.RecordScrapeDuration(float64(time.Since(start).Nanoseconds()))
}

func (s *Scraper) updateStats(resp *client.StatsResponse) {
	metrics.UpdateStats(
		resp.TotalPromptTokens,
		resp.TotalCompletionTokens,
		resp.TotalCachedTokens,
		resp.TotalRequests,
	)
	metrics.UpdateGauges(
		resp.CacheEfficiency,
		resp.AvgPrefillTPS,
		resp.AvgGenerationTPS,
		resp.UptimeSeconds,
	)
}

func (s *Scraper) updateActiveModels(resp *client.StatsResponse) {
	am := resp.ActiveModels
	metrics.UpdateActiveRequests(am.TotalActiveRequests, am.TotalWaitingRequests)
	metrics.UpdateMemoryPressure(
		am.MemoryPressure.PressureLevel,
		am.MemoryPressure.CurrentBytes,
		am.MemoryPressure.SoftBytes,
		am.MemoryPressure.HardBytes,
	)

	metrics.ResetModelGauges()
	metrics.ResetRequestGauges()
	metrics.ResetWaitingGauges()

	for _, m := range am.Models {
		model := m.ModelID

		metrics.SetModelWaitingRequests(model, float64(len(m.Waiting)))
		metrics.SetModelPrefillingRequests(model, float64(len(m.Prefilling)))
		metrics.SetModelGeneratingRequests(model, float64(len(m.Generating)))
		metrics.SetModelActiveRequests(model, float64(m.ActiveRequests))

		for _, g := range m.Generating {
			metrics.SetModelGeneratingTPS(model, g.TokensPerSecond)
			metrics.SetModelGeneratingTokens(model, float64(g.GeneratedTokens))
			metrics.SetModelGeneratingElapsed(model, g.ElapsedSeconds)
		}

		for _, w := range m.Waiting {
			metrics.SetModelWaitingQueuePos(model, float64(w.QueuePosition))
			metrics.SetModelWaitingElapsed(model, w.ElapsedSeconds)
		}
	}
}
