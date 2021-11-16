# CloudFlare Exporter

## Description
The project drives metrics from the [CloudFlare graphQL API](https://developers.cloudflare.com/analytics/graphql-api) to [Prometheus](https://prometheus.io/).

## CloudFlare Data Extraction
CloudFlare graphQL API updates data every minute.
To keep the system compliant with Prometheus logic, the exporter will aggregate data in the time window from
_four minutes back_ to _three minutes back_.

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

There are two log levels, `DEBUG` and `INFO`. The `DEBUG` log level is the default.
Set `ENV` variable to `prod` to run with `INFO` log level.

## Variables
The `SUBSYSTEM` environment variable will be used as [metrics subsystem](https://github.com/prometheus/client_golang/blob/master/prometheus/examples_test.go#L38). The variable can be used to differentiate metric names for different CloudFlare domain zones.
The default value is `gather_town`.

## Metrics Exposed
Default metric names:

| name |
|-----------------------------------------------------|
| cloudflare_gather_town_browser_map_page_views_count |
| cloudflare_gather_town_country_map_bytes_sum        |
| cloudflare_gather_town_country_map_requests_count   |
| cloudflare_gather_town_country_map_threats_count    |
| cloudflare_gather_town_response_bytes_count         |
| cloudflare_gather_town_response_status_count        |
| cloudflare_gather_town_visits_count                 |

For detailed explanation look at the metric description in the source code, Grafana or Prometheus.

# ROADMAP
* Add metrics from [CloudFlare browser insights](https://support.cloudflare.com/hc/en-us/articles/360033929991-Cloudflare-Browser-Insights)


## License
See [Licence](LICENSE)
