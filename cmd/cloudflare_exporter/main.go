package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	config "github.com/gathertown/cloudflare_exporter/internal/config"
	r "github.com/gathertown/cloudflare_exporter/pkg/cloudflare/requests"
	t "github.com/gathertown/cloudflare_exporter/pkg/cloudflare/traffic"
	log "github.com/gathertown/cloudflare_exporter/pkg/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	metrics "github.com/gathertown/cloudflare_exporter/internal/metrics"
)

const (
	rfc3339 = "2006-01-02T15:04:05-0700"
)

var qr = r.Q
var qt = t.Q
var cfg = config.FromEnv()
var logger = log.New(os.Stdout, cfg.Env)

// timeWindow will return the exact minute in time.RFC3339
// e.g. start: 2021-05-06T09:55:00Z, end: 2021-05-06T09:56:00Z
func timeWindow(start string, end string) (string, string, error) {
	t := time.Now().UTC() //Get the current time
	f := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:00-0000",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute())
	c, err := time.Parse(rfc3339, f)

	d1, err := time.ParseDuration(start)
	if err != nil {
		return "", "", err
	}

	d2, err := time.ParseDuration(end)
	if err != nil {
		return "", "", err
	}

	t1 := c.UTC().Add(d1).Format(time.RFC3339)
	t2 := c.UTC().Add(d2).Format(time.RFC3339)
	return t1, t2, nil
}

func recordMetrics() {
	go func() {
		for {
			currentMinute := time.Now().UTC().Minute()
			t1, t2, err := timeWindow("-4m", "-3m")
			rData, err := r.Run(cfg.ZoneTag, cfg.Token, t1, t2)
			tData, err := t.Run(cfg.ZoneTag, cfg.Token, t1, t2, 1000)
			if err != nil {
				panic(err)
			}

			for _, k := range tData {
				metrics.PoolHealthStatus.WithLabelValues(k.ColoCode, k.LbName, k.OriginName, k.Policy, k.PoolName, k.Region).Set(float64(k.Healthy))
			}

			for _, k := range rData {
				for _, v := range k.HttpRequestsAdaptiveGroups {
					v1, err := strconv.ParseFloat(fmt.Sprintf("%v", v.Sum.Visits), 64)
					if err != nil {
						logger.Info("conversion error", "error", err)
					}
					metrics.EdgeVisits.Add(v1)

					v2, err := strconv.ParseFloat(fmt.Sprintf("%v", v.Sum.EdgeResponseBytes), 64)
					if err != nil {
						logger.Info("conversion error", "error", err)
					}
					metrics.EdgeBytes.Add(v2)
				}

				// Extract HttpRequests1mGroups data
				for _, v := range k.HttpRequests1mGroups {
					for _, d := range v.Sum.ResponseStatusMap {
						v1, err := strconv.ParseFloat(fmt.Sprintf("%v", d.Requests), 64)
						if err != nil {
							logger.Info("conversion error", "error", err)
						}
						metrics.EdgeResponseStatus.WithLabelValues(fmt.Sprintf("%v", d.EdgeResponseStatus)).Add(v1)
					}
					for _, b := range v.Sum.BrowserMap {
						v1, err := strconv.ParseFloat(fmt.Sprintf("%v", b.PageViews), 64)
						if err != nil {
							logger.Info("conversion error", "error", err)
						}
						metrics.EdgeBrowserMap.WithLabelValues(fmt.Sprintf("%v", b.UaBrowserFamily)).Add(v1)
					}
					for _, c := range v.Sum.CountryMap {
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
						metrics.EdgeCountryMapRequests.WithLabelValues(fmt.Sprintf("%v", c.ClientCountryName)).Add(v1)
						metrics.EdgeCountryMapBytes.WithLabelValues(fmt.Sprintf("%v", c.ClientCountryName)).Add(v2)
						metrics.EdgeCountryMapThreats.WithLabelValues(fmt.Sprintf("%v", c.ClientCountryName)).Add(v3)
					}
				}
			}

			// Don't expose twice the same data. If we already reported this minute, wait for the next one.
			for {
				if time.Now().UTC().Minute() != currentMinute {
					logger.Debug("Updating data", "currentMinute", currentMinute, "newMinute", time.Now().UTC().Minute())
					break
				}
				logger.Debug("Sleeping 35 seconds", "currentMinute", currentMinute)
				time.Sleep(35 * time.Second)
			}
		}
	}()
}

// run labels nodes if label is missing
func main() {
	cfg := config.FromEnv()
	recordMetrics()
	logger := log.New(os.Stdout, cfg.Env)
	logger.Info("Launching cloudformation_exporter", "zoneTag", cfg.ZoneTag, "metrics subsystem", cfg.Sub)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
