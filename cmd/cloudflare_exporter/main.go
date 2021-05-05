package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	config "github.com/gathertown/cloudflare_exporter/internal/config"
	common "github.com/gathertown/cloudflare_exporter/pkg"
	r "github.com/gathertown/cloudflare_exporter/pkg/cloudflare/requests"
	log "github.com/gathertown/cloudflare_exporter/pkg/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var q = common.Q
var cfg = config.FromEnv()
var logger = log.New(os.Stdout, cfg.Env)

// define custom metrics
// https://pkg.go.dev/github.com/prometheus/client_golang@v1.10.0/prometheus#GaugeVec
var (
	edgeVisits = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cloudflare_exporter_visits_sum",
		Help: "Sum of processed events",
	})

	edgeBytes = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cloudflare_exporter_edge_response_bytes_sum",
		Help: "Sum of response bytes",
	})

	edgeBrowserMap = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cloudflare_exporter_edge_browser_map_page_views_sum",
			Help: "Sum of page views per browser",
		},
		[]string{"family"},
	)

	edgeCountryMapRequests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cloudflare_exporter_country_map_requests_sum",
			Help: "Sum of requests per country",
		},
		[]string{"country"},
	)

	edgeCountryMapBytes = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cloudflare_exporter_country_map_bytes_sum",
			Help: "Sum of bytes per country",
		},
		[]string{"country"},
	)

	edgeCountryMapThreats = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cloudflare_exporter_country_map_threats_sum",
			Help: "Sum of threats per country",
		},
		[]string{"country"},
	)

	edgeResponseStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cloudflare_exporter_response_status_sum",
			Help: "Sum of responses per status code",
		},
		[]string{"status"},
	)
)

func recordMetrics() {
	go func() {
		for {
			data, err := r.Requests(cfg.ZoneTag, cfg.Token)
			if err != nil {
				panic(err)
			}

			for _, k := range data {
				// Extract HttpRequestsAdaptiveGroups data
				for _, v := range k.HttpRequestsAdaptiveGroups {
					logger.Debug("HttpRequestsAdaptiveGroups", "datetime", v.Dimensions.DatetimeMinute, "visits", v.Sum.Visits, "bytes", v.Sum.EdgeResponseBytes)
					v1, err := strconv.ParseFloat(fmt.Sprintf("%v", v.Sum.Visits), 64)
					if err != nil {
						logger.Info("conversion error", "error", err)
					}
					edgeVisits.Set(v1)

					v2, err := strconv.ParseFloat(fmt.Sprintf("%v", v.Sum.EdgeResponseBytes), 64)
					if err != nil {
						logger.Info("conversion error", "error", err)
					}
					edgeBytes.Set(v2)
				}

				// Extract HttpRequests1mGroups data
				for _, v := range k.HttpRequests1mGroups {
					for _, d := range v.Sum.ResponseStatusMap {
						logger.Debug("responseStatus", "status", d.EdgeResponseStatus, "requests", d.Requests)

						v1, err := strconv.ParseFloat(fmt.Sprintf("%v", d.Requests), 64)
						if err != nil {
							logger.Info("conversion error", "error", err)
						}
						edgeBrowserMap.WithLabelValues(fmt.Sprintf("%v", d.EdgeResponseStatus)).Set(v1)
					}
					for _, b := range v.Sum.BrowserMap {
						logger.Debug("BrowserMap", "browser", b.UaBrowserFamily, "views", b.PageViews)

						v1, err := strconv.ParseFloat(fmt.Sprintf("%v", b.PageViews), 64)
						if err != nil {
							logger.Info("conversion error", "error", err)
						}
						edgeBrowserMap.WithLabelValues(fmt.Sprintf("%v", b.UaBrowserFamily)).Set(v1)
					}
					for _, c := range v.Sum.CountryMap {
						logger.Debug("CountryMap", "countryCode", c.ClientCountryName, "bytes", c.Bytes, "requests", c.Requests, "threats", c.Threats)

						v1, err := strconv.ParseFloat(fmt.Sprintf("%v", c.Requests), 64)
						if err != nil {
							logger.Info("conversion error", "error", err)
						}

						v2, err := strconv.ParseFloat(fmt.Sprintf("%v", c.Bytes), 64)
						if err != nil {
							logger.Info("conversion error", "error", err)
						}

						v3, err := strconv.ParseFloat(fmt.Sprintf("%v", c.Threats), 64)
						if err != nil {
							logger.Info("conversion error", "error", err)
						}

						// set metrics
						edgeCountryMapRequests.WithLabelValues(fmt.Sprintf("%v", c.ClientCountryName)).Set(v1)
						edgeCountryMapBytes.WithLabelValues(fmt.Sprintf("%v", c.ClientCountryName)).Set(v2)
						edgeCountryMapThreats.WithLabelValues(fmt.Sprintf("%v", c.ClientCountryName)).Set(v3)
					}
				}
			}
			// TODO: Use an accurate method to fetch data. We can pass using "time.Now()"
			// We're assuming the call and data processing will take ~2s as the graphQL is not
			// extremely responsive
			time.Sleep(58 * time.Second)
		}
	}()
}

// run labels nodes if label is missing
func main() {
	cfg := config.FromEnv()
	recordMetrics()
	logger := log.New(os.Stdout, cfg.Env)
	logger.Info("Launching cloudformation_exporter", "zoneTag", cfg.ZoneTag)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
