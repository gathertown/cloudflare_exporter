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
		Name:      "external_visits_count",
		Namespace: namespace,
		Subsystem: subsystem,
		Help:      "Counts the number of requests by end-users that were initiated from a different website (i.e. where the request HTTP Referer header does not match the host in the HTTP Host header)",
	})

	EdgeBytes = promauto.NewCounter(prometheus.CounterOpts{
		Name:      "response_bytes_count",
		Namespace: namespace,
		Subsystem: subsystem,
		Help:      "Counts the amount of data transferred from Cloudflare to end users within a certain period of time. Total bandwidth equals the sum of all EdgeResponseBytes for a certain period of time",
	})

	EdgeBrowserMap = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "browser_map_page_views_count",
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      "Counts the successful requests for HTML",
		},
		[]string{"family"},
	)

	EdgeCountryMapRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "country_map_requests_count",
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      "Counts the country from which request originated",
		},
		[]string{"country"},
	)

	EdgeCountryMapBytes = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "country_map_bytes_count",
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      "Count of bytes returned to client per country",
		},
		[]string{"country"},
	)

	EdgeCountryMapThreats = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "country_map_threats_count",
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      "Counts requests classified as threats per country",
		},
		[]string{"country"},
	)

	EdgeResponseStatus = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "response_status_count",
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      "Counts HTTP response status code returned to client",
		},
		[]string{"status"},
	)
)
