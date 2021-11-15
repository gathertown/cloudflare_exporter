package metrics

import (
	config "github.com/gathertown/cloudflare_exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Configure global package variables
var cfg = config.FromEnv()
var namespace = "cloudflare"
var subsystem = cfg.Sub

// define custom metrics
// https://pkg.go.dev/github.com/prometheus/client_golang@v1.10.0/prometheus#GaugeVec
var (
	EdgeVisits = promauto.NewCounter(prometheus.CounterOpts{
		Name:      "visits_count",
		Namespace: namespace,
		Subsystem: subsystem,
		Help:      "Count of visits",
	})

	EdgeBytes = promauto.NewGauge(prometheus.GaugeOpts{
		Name:      "response_bytes_sum",
		Namespace: namespace,
		Subsystem: subsystem,
		Help:      "Sum of response bytes",
	})

	EdgeBrowserMap = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "browser_map_page_views_count",
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      "Count of page views per browser",
		},
		[]string{"family"},
	)

	EdgeCountryMapRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "country_map_requests_count",
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      "Count of requests per country",
		},
		[]string{"country"},
	)

	EdgeCountryMapBytes = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "country_map_bytes_count",
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      "Count of bytes per country",
		},
		[]string{"country"},
	)

	EdgeCountryMapThreats = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "country_map_threats_count",
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      "Count of threats per country",
		},
		[]string{"country"},
	)

	EdgeResponseStatus = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "response_status_count",
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      "Count of responses per status code",
		},
		[]string{"status"},
	)
)
