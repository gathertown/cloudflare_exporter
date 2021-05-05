package common

import graphql "github.com/shurcooL/graphql"

// -start graphQL structs

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

type Zones []struct {
	HttpRequests1mGroups       httpRequests1mGroups       `graphql:"httpRequests1mGroups(filter: {datetime_geq: $dateTimeStart, datetime_leq: $dateTimeEnd}, orderBy: [ datetimeMinute_ASC], limit: 1)"`
	HttpRequestsAdaptiveGroups httpRequestsAdaptiveGroups `graphql:"httpRequestsAdaptiveGroups(filter: {datetime_geq: $dateTimeStart, datetime_leq: $dateTimeEnd}, orderBy: [ datetimeMinute_ASC], limit: 1)"`
}

var Q struct {
	Viewer struct {
		Zones Zones `graphql:"zones(filter: {zoneTag: $zoneTag})"`
	}
}

// -end graphQL structs
