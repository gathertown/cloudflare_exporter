# CloudFlare Exporter

## Description

The project drives metrics from the [CloudFlare graphQL API](https://developers.cloudflare.com/analytics/graphql-api) to [Prometheus](https://prometheus.io/).

The API exposes sum of metrics for a predefined time window. The time window is approximately one minute.

Metrics need to be aggregated on CloudFlare. There are no metrics exposed for the _last minute_. The time window starts _three minutes_ back and ends _two minutes_ back.

## Develop

Run the app:

```bash
export TOKEN="<cloudflare_token>"
export ZONETAG="<cloudflare_zone_tag>"
export SUBSYSTEM="<my_org>"
make run
```

Build the app:

```bash
make build
./bin/cloudflare_exporter
```

## Metrics

The `SUBSYSTEM` environment variable will be used as [metrics subsystem](https://github.com/prometheus/client_golang/blob/master/prometheus/examples_test.go#L38).
Can be used to differentiate metric names for different cloudflare domain zones.

The default `SUBSYSTEM` is `gather_town`. Default metric names:

| name |
|----------------------------------------------------------------------------|
|cloudflare_gather_town_browser_map_page_views_sum{family="<browser_family>"}|
|cloudflare_gather_town_country_map_bytes_sum{country="<country>"}           |
|cloudflare_gather_town_country_map_requests_sum{country="<country>"}        |
|cloudflare_gather_town_country_map_threats_sum{country="<country>"}         |
|cloudflare_gather_town_response_bytes_sum                                   |
|cloudflare_gather_town_response_status_sum{status="<status_code>"}          |
|cloudflare_gather_town_visits_sum                                           |

All metrics are [gauges](https://prometheus.io/docs/concepts/metric_types/#gauge).

## License
See [Licence](LICENSE)
