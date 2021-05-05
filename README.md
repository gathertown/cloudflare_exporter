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
make run
```

Build the app:

```bash
make build
./bin/cloudflare_exporter
```





## Metrics

| name |
|----------------------------------------------------------------------------|
|gather_town_cloudflare_browser_map_page_views_sum{family="<browser_family>"}|
|gather_town_cloudflare_country_map_bytes_sum{country="<country>"}           |
|gather_town_cloudflare_country_map_requests_sum{country="<country>"}        |
|gather_town_cloudflare_country_map_threats_sum{country="<country>"}         |
|gather_town_cloudflare_response_bytes_sum                                   |
|gather_town_cloudflare_response_status_sum{status="<status_code>"}          |

All metrics are [gauges](https://prometheus.io/docs/concepts/metric_types/#gauge).

## License
See [Licence](LICENSE)