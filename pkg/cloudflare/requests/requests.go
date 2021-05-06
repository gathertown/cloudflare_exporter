package requests

import (
	"context"
	"os"

	"github.com/gathertown/cloudflare_exporter/internal/config"
	common "github.com/gathertown/cloudflare_exporter/pkg"
	log "github.com/gathertown/cloudflare_exporter/pkg/log"
	graphql "github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

var cfg = config.FromEnv()
var logger = log.New(os.Stdout, cfg.Env)
var q = common.Q

type z = common.Zones

const endpoint = "https://api.cloudflare.com/client/v4/graphql"

func Requests(zoneTag string, token string, start string, end string) (z, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := graphql.NewClient(endpoint, httpClient)

	// Start duration
	// d1, err := time.ParseDuration("-3m")
	// if err != nil {
	// 	return nil, err
	// }

	// Stop duration
	// d2, err := time.ParseDuration("-2m")
	// if err != nil {
	// 	return nil, err
	// }

	// set start and stop time
	// t1 := time.Now().UTC().Add(d1).Format(time.RFC3339)
	// t2 := time.Now().UTC().Add(d2).Format(time.RFC3339)

	// set request options
	opts := map[string]interface{}{
		"dateTimeStart": graphql.String(start),
		"dateTimeEnd":   graphql.String(end),
		"zoneTag":       graphql.String(zoneTag),
	}

	// perform graphQL query
	err := client.Query(context.Background(), &q, opts)
	if err != nil {
		return nil, err
	}

	// logger.Debug("time", "dateTimeStart", t1)
	// logger.Debug("time", "dateTimeEnd", t2)

	// // keep loop for debugging purposes
	// for _, k := range q.Viewer.Zones {
	// 	for _, v := range k.HttpRequestsAdaptiveGroups {
	// 		logger.Debug("HttpRequestsAdaptiveGroups", "datetime", v.Dimensions.DatetimeMinute, "visits", v.Sum.Visits, "bytes", v.Sum.EdgeResponseBytes)
	// 	}
	// 	for _, v := range k.HttpRequests1mGroups {
	// 		fmt.Printf("DateTime: %s\n", v.Dimensions.DatetimeMinute)
	// 		for _, d := range v.Sum.ResponseStatusMap {
	// 			logger.Debug("responseStatus", "status", d.EdgeResponseStatus, "requests", d.Requests)
	// 		}
	// 		for _, b := range v.Sum.BrowserMap {
	// 			logger.Debug("BrowserMap", "browser", b.UaBrowserFamily, "views", b.PageViews)
	// 		}
	// 		for _, c := range v.Sum.CountryMap {
	// 			logger.Debug("CountryMap", "countryCode", c.ClientCountryName, "bytes", c.Bytes, "requests", c.Requests, "threats", c.Threats)
	// 		}
	// 	}
	// }

	return q.Viewer.Zones, nil
}
