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
	NowPlayingSubmits   *prometheus.CounterVec
	ScrobbleSubmits     *prometheus.CounterVec
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
	nowPlayingSubmits := promauto.With(reg).NewCounterVec(prometheus.CounterOpts{Name: "scrobblet_now_playing_submits_total", Help: "Total number of now playing submits."}, []string{"target", "status"})
	scrobbleSubmits := promauto.With(reg).NewCounterVec(prometheus.CounterOpts{Name: "scrobblet_scrobble_submits_total", Help: "Total number of scrobble submits."}, []string{"target", "status"})
	scrobbleEvaluations := promauto.With(reg).NewCounterVec(prometheus.CounterOpts{Name: "scrobblet_scrobble_evaluations_total", Help: "Total number of track evaluations."}, []string{"result"})
	pollDuration := promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{Name: "scrobblet_poll_duration_seconds", Help: "Duration (in seconds) of polls.", Buckets: prometheus.DefBuckets}, []string{"source"})
	submitDuration := promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{Name: "scrobblet_submit_duration_seconds", Help: "Duration (in seconds) of submits.", Buckets: prometheus.DefBuckets}, []string{"target", "kind"})

	return &Metrics{
		Registry:            reg,
		BuildInfo:           buildInfo,
		PollingInterval:     pollingInterval,
		Polls:               polls,
		NowPlayingSubmits:   nowPlayingSubmits,
		ScrobbleSubmits:     scrobbleSubmits,
		ScrobbleEvaluations: scrobbleEvaluations,
		PollDuration:        pollDuration,
		SubmitDuration:      submitDuration,
	}
}
