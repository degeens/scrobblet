package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	StatusSuccess  = "success"
	StatusFailure  = "failure"
	Met            = "threshold_met"
	NotMet         = "threshold_not_met"
	KindNowPlaying = "now_playing"
	KindScrobble   = "scrobble"
)

type Metrics struct {
	Registry            *prometheus.Registry
	BuildInfo           *prometheus.GaugeVec
	PollingInterval     prometheus.Gauge
	Polls               *prometheus.CounterVec
	Submits             *prometheus.CounterVec
	ScrobbleEvaluations *prometheus.CounterVec
	PollDuration        *prometheus.HistogramVec
	SubmitDuration      *prometheus.HistogramVec
}

func New() *Metrics {
	reg := prometheus.NewRegistry()

	reg.MustRegister(collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	buildInfo := promauto.With(reg).NewGaugeVec(prometheus.GaugeOpts{Name: "scrobblet_build_info", Help: "Build information."}, []string{"version"})
	pollingInterval := promauto.With(reg).NewGauge(prometheus.GaugeOpts{Name: "scrobblet_polling_interval_seconds", Help: "Current polling interval (in seconds)."})
	polls := promauto.With(reg).NewCounterVec(prometheus.CounterOpts{Name: "scrobblet_polls_total", Help: "Total number of polls."}, []string{"source", "status"})
	submits := promauto.With(reg).NewCounterVec(prometheus.CounterOpts{Name: "scrobblet_submits_total", Help: "Total number of submits."}, []string{"target", "kind", "status"})
	scrobbleEvaluations := promauto.With(reg).NewCounterVec(prometheus.CounterOpts{Name: "scrobblet_scrobble_evaluations_total", Help: "Total number of scrobble evaluations."}, []string{"result"})
	pollDuration := promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{Name: "scrobblet_poll_duration_seconds", Help: "Duration (in seconds) of polls.", Buckets: prometheus.DefBuckets}, []string{"source"})
	submitDuration := promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{Name: "scrobblet_submit_duration_seconds", Help: "Duration (in seconds) of submits.", Buckets: prometheus.DefBuckets}, []string{"target", "kind"})

	return &Metrics{
		Registry:            reg,
		BuildInfo:           buildInfo,
		PollingInterval:     pollingInterval,
		Polls:               polls,
		Submits:             submits,
		ScrobbleEvaluations: scrobbleEvaluations,
		PollDuration:        pollDuration,
		SubmitDuration:      submitDuration,
	}
}
