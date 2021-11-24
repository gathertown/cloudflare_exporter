package requests

import (
	"context"
	"os"

	"github.com/gathertown/cloudflare_exporter/internal/config"
	log "github.com/gathertown/cloudflare_exporter/pkg/log"
	graphql "github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

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
// Number "unique visitors" per colocation is not supported in httpRequestsAdaptiveGroups, but the httpRequestsAdaptiveGroups API does support visits.
// A visit is defined as a page view that originated from a different website or direct link. Cloudflare checks where the HTTP referer does not match the hostname. One visit can consist of multiple page views.
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

type Zones []struct {
	HttpRequests1mGroups       httpRequests1mGroups       `graphql:"httpRequests1mGroups(filter: {datetime_geq: $dateTimeStart, datetime_leq: $dateTimeEnd}, orderBy: [ datetimeMinute_ASC], limit: 1)"`
	HttpRequestsAdaptiveGroups httpRequestsAdaptiveGroups `graphql:"httpRequestsAdaptiveGroups(filter: {datetime_geq: $dateTimeStart, datetime_leq: $dateTimeEnd}, orderBy: [ datetimeMinute_ASC], limit: 1)"`
}

var Q struct {
	Viewer struct {
		Zones Zones `graphql:"zones(filter: {zoneTag: $zoneTag})"`
	}
}

var cfg = config.FromEnv()
var logger = log.New(os.Stdout, cfg.Env)

const endpoint = "https://api.cloudflare.com/client/v4/graphql"

// Run executes the graphQL query
func Run(zoneTag string, token string, start string, end string) (Zones, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := graphql.NewClient(endpoint, httpClient)

	// set request options
	opts := map[string]interface{}{
		"dateTimeStart": graphql.String(start),
		"dateTimeEnd":   graphql.String(end),
		"zoneTag":       graphql.String(zoneTag),
	}

	// perform graphQL query
	err := client.Query(context.Background(), &Q, opts)
	if err != nil {
		return nil, err
	}

	return Q.Viewer.Zones, nil
}
