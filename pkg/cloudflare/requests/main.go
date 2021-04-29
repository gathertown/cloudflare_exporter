package requests

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gathertown/cloudflare_exporter/internal/config"
	log "github.com/gathertown/cloudflare_exporter/pkg/log"
	graphql "github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

var cfg = config.FromEnv()
var logger = log.New(os.Stdout, cfg.Env)

const endpoint = "https://api.cloudflare.com/client/v4/graphql"

// start graphQL structs

type countryMap []struct {
	ClientCountryName graphql.String
	Requests          graphql.Int
	Threats           graphql.Int
	Bytes             graphql.Int
}

type browserMap []struct {
	UaBrowserFamily graphql.String
	PageViews       graphql.Int
}

type responseStatusMap []struct {
	Requests           graphql.Int
	EdgeResponseStatus graphql.Int
}

type httpRequests1mGroupsSum struct {
	CountryMap        countryMap
	BrowserMap        browserMap
	ResponseStatusMap responseStatusMap
}

type dimensions struct {
	DatetimeMinute graphql.String
}

type httpRequests1mGroups []struct {
	Dimensions dimensions
	Sum        httpRequests1mGroupsSum
}

// NOTE:
// "unique visitors" per colocation is not supported in httpRequestsAdaptiveGroups,
// but the httpRequestsAdaptiveGroups API does support visits.
// A visit is defined as a page view that originated from a different website or direct link.
// Cloudflare checks where the HTTP referer does not match the hostname. One visit can consist of multiple page views.
//
// URL: https://developers.cloudflare.com/analytics/graphql-api/migration-guides/graphql-api-analytics
type httpRequestsAdaptiveGroupsSum struct {
	Visits            graphql.Int
	EdgeResponseBytes graphql.Int
}

type httpRequestsAdaptiveGroups []struct {
	Dimensions dimensions
	Sum        httpRequestsAdaptiveGroupsSum
}

type zones []struct {
	HttpRequests1mGroups       httpRequests1mGroups       `graphql:"httpRequests1mGroups(filter: {datetime_geq: $dateTimeStart, datetime_leq: $dateTimeEnd}, orderBy: [ datetimeMinute_ASC], limit: 1)"`
	HttpRequestsAdaptiveGroups httpRequestsAdaptiveGroups `graphql:"httpRequestsAdaptiveGroups(filter: {datetime_geq: $dateTimeStart, datetime_leq: $dateTimeEnd}, orderBy: [ datetimeMinute_ASC], limit: 1)"`
}

var q struct {
	Viewer struct {
		Zones zones `graphql:"zones(filter: {zoneTag: $zoneTag})"`
	}
}

func Requests(zoneTag string, token string) (int, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := graphql.NewClient(endpoint, httpClient)

	// Start duration
	d1, err := time.ParseDuration("-3m")
	if err != nil {
		return 0, err
	}

	// Stop duration
	d2, err := time.ParseDuration("-2m")
	if err != nil {
		return 0, err
	}
	t1 := time.Now().UTC().Add(d1).Format(time.RFC3339)
	t2 := time.Now().UTC().Add(d2).Format(time.RFC3339)

	opts := map[string]interface{}{
		"dateTimeStart": graphql.String(t1),
		"dateTimeEnd":   graphql.String(t2),
		"zoneTag":       graphql.String(zoneTag),
	}

	err = client.Query(context.Background(), &q, opts)
	if err != nil {
		return 0, err
	}

	logger.Debug("time", "dateTimeStart", t1)
	logger.Debug("time", "dateTimeEnd", t2)

	for _, k := range q.Viewer.Zones {
		for _, v := range k.HttpRequestsAdaptiveGroups {
			fmt.Printf("DateTime: %s, Visits: %v, Bytes: %v\n", v.Dimensions.DatetimeMinute, v.Sum.Visits, v.Sum.EdgeResponseBytes)
		}
		for _, v := range k.HttpRequests1mGroups {
			fmt.Printf("DateTime: %s\n", v.Dimensions.DatetimeMinute)
			for _, d := range v.Sum.ResponseStatusMap {
				fmt.Printf("status code: %v, requests: %v\n", d.EdgeResponseStatus, d.Requests)
			}
			for _, b := range v.Sum.BrowserMap {
				fmt.Printf("BrowserFamily: %v, PageViews: %v\n", b.UaBrowserFamily, b.PageViews)
			}
			for _, c := range v.Sum.CountryMap {
				fmt.Printf("Country: %v, Bytes: %v, Requests: %v, Threats: %v\n", c.ClientCountryName, c.Bytes, c.Requests, c.Threats)
			}
		}
	}
	return 1, nil
}
