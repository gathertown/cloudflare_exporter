package traffic

import (
	"context"
	"os"
	"strings"

	"github.com/gathertown/cloudflare_exporter/internal/config"
	log "github.com/gathertown/cloudflare_exporter/pkg/log"
	graphql "github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

// -start graphQL structs
type dimensions struct {
	ColoCode graphql.String
	LbName   graphql.String
	Region   graphql.String
	// SelectedOriginIndex    graphql.Int
	SelectedOriginName graphql.String
	// SelectedPoolAvgRttMs   graphql.Float
	SelectedPoolHealthy graphql.Int
	// SelectedPoolId         graphql.String
	// SelectedPoolIndex      graphql.Int
	SelectedPoolName graphql.String
	// SessionAffinityEnabled graphql.String
	SteeringPolicy graphql.String
}

type loadBalancingRequestsAdaptiveGroups []struct {
	Dimensions dimensions
}

type Zones []struct {
	LoadBalancingRequestsAdaptiveGroups loadBalancingRequestsAdaptiveGroups `graphql:"loadBalancingRequestsAdaptiveGroups(filter: {datetime_geq: $dateTimeStart, datetime_lt: $dateTimeEnd, selectedPoolHealthy: $health, selectedPoolId: $poolID}, limit: $limit)"`
}

var Q struct {
	Viewer struct {
		Zones Zones `graphql:"zones(filter: {zoneTag: $zoneTag})"`
	}
}

// -end graphQL structs

type Pool struct {
	ColoCode   string
	Healthy    int
	LbName     string
	OriginName string
	Policy     string
	PoolName   string
	Region     string
}

var cfg = config.FromEnv()
var logger = log.New(os.Stdout, cfg.Env)
var origins = strings.Split(cfg.Origin, ",")

const endpoint = "https://api.cloudflare.com/client/v4/graphql"

// getOriginIDs returns a slice of oririns
// func getOriginIDs(token string, loadBalancer string) ([]string, error) {
// 	origins := []string{}
// 	api, err := cloudflare.NewWithAPIToken(token)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Fetch the zone ID for zone example.org
// 	zoneID, err := api.ZoneIDByName(loadBalancer)
// 	ctx := context.Background()
// 	if err != nil {
// 		return nil, err
// 	}

// 	lbs, err := api.ListLoadBalancers(ctx, zoneID)
// 	for _, v := range lbs {
// 		for _, p := range v.DefaultPools {
// 			origins = append(origins, p)
// 		}
// 	}
// 	return origins, nil
// }

// fetchData returns a slice of structs containing metrics info
func fetchData(zoneTag string, token string, start string, end string, limit int) ([]Pool, error) {

	p := []Pool{}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := graphql.NewClient(endpoint, httpClient)

	// set request options
	opts1 := map[string]interface{}{
		"dateTimeStart": graphql.String(start),
		"dateTimeEnd":   graphql.String(end),
		"zoneTag":       graphql.String(zoneTag),
		"health":        graphql.Int(1),
		"limit":         graphql.Int(limit),
	}

	opts2 := map[string]interface{}{
		"dateTimeStart": graphql.String(start),
		"dateTimeEnd":   graphql.String(end),
		"zoneTag":       graphql.String(zoneTag),
		"health":        graphql.Int(0),
		"limit":         graphql.Int(limit),
	}

	err := client.Query(context.Background(), &Q, opts1)
	if err != nil {
		return p, err
	}

	// return Q.Viewer.Zones, nil
	for _, k := range Q.Viewer.Zones {
		for _, v := range k.LoadBalancingRequestsAdaptiveGroups {
			s1 := Pool{
				ColoCode:   string(v.Dimensions.ColoCode),
				Healthy:    int(v.Dimensions.SelectedPoolHealthy),
				LbName:     string(v.Dimensions.LbName),
				OriginName: string(v.Dimensions.SelectedOriginName),
				Policy:     string(v.Dimensions.SteeringPolicy),
				PoolName:   string(v.Dimensions.SelectedPoolName),
				Region:     string(v.Dimensions.Region),
			}
			// Add entry to pools if missing
			if len(p) > 1 {
				found := false
				for _, v1 := range p {
					if v1.OriginName == string(v.Dimensions.SelectedOriginName) {
						found = true
					}
				}
				if !found {
					p = append(p, s1)
				}
			} else {
				p = append(p, s1)
			}
		}
	}

	err = client.Query(context.Background(), &Q, opts2)
	if err != nil {
		return nil, err
	}

	for _, k := range Q.Viewer.Zones {
		if len(k.LoadBalancingRequestsAdaptiveGroups) > 1 {
			for _, v := range k.LoadBalancingRequestsAdaptiveGroups {
				s1 := Pool{
					ColoCode:   string(v.Dimensions.ColoCode),
					Healthy:    int(v.Dimensions.SelectedPoolHealthy),
					LbName:     string(v.Dimensions.LbName),
					OriginName: string(v.Dimensions.SelectedOriginName),
					Policy:     string(v.Dimensions.SteeringPolicy),
					PoolName:   string(v.Dimensions.SelectedPoolName),
					Region:     string(v.Dimensions.Region),
				}
				// Add entry to pools if missing
				if len(p) > 1 {
					for _, v1 := range p {
						if v1.OriginName == string(v.Dimensions.SelectedOriginName) {
							logger.Info("Found unhealthy pools", "poolName", v1.OriginName)
							v1.Healthy = int(v.Dimensions.SelectedPoolHealthy)
						}
					}
				} else {
					logger.Info("Found unhealthy pools", "poolName", string(v.Dimensions.SelectedOriginName))
					p = append(p, s1)
				}
			}
		}
	}

	return p, nil
}

// Run executes the graphQL query to fetch traffic (load balancer) related statistics
func Run(zoneTag string, token string, start string, end string, limit int) ([]Pool, error) {
	data, err := fetchData(zoneTag, token, start, end, limit)
	if err != nil {
		return nil, err
	}

	return data, nil
}
