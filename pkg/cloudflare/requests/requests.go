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

	return q.Viewer.Zones, nil
}
