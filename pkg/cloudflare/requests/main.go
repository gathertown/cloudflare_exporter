// Github Crawler golang package fetches list of repos featuring a topic, returns a map

package requests

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// Crawler struct defines the fields needed by ghcrawler pkg
type Crawler struct {
	PageItems int    // The number of items to fetch per API call, the limit is 100
	Org       string // Github organisation
	Token     string // Github token
	Topic     string // Github topic
}

/* GraphQL struct start */
type pageinfo struct {
	HasNextPage bool
	EndCursor   string // When 'HasNextPage' is 'true', EndCursor will return the page number to fetch next
}

type branch struct {
	Name string
}

type repository struct {
	DefaultBranchRef branch
	Name             string // Github repository name
	URL              string // Github repository URL
	ID               string // Github uses an alphanumeric string as a repository unique identifier
}

type repos []struct {
	Repository repository `graphql:"... on Repository"`
}

type search struct {
	PageInfo pageinfo // Contains info about graphQL API pagination
	Repos    repos    `graphql:"repo: nodes"`
}

func client(ctx context.Context, token string) *githubv4.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return githubv4.NewClient(tc)
}

/* GraphQL struct end */

// Fetch gathers repositories from github featuring a topic. Returns the names of the repositories as a slice of strings.
func (c Crawler) Fetch(ctx context.Context) ([]map[string]string, error) {
	// Github API v4 requires an API Token:
	// https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line
	clnt := client(ctx, c.Token)

	// GraphQL does not support string concatenation. We'll provide the query as a pre-constructed string instead.
	// https://github.com/shurcooL/githubv4/issues/19#issuecomment-326034518
	query := fmt.Sprintf("topic:%s org:%s", c.Topic, c.Org)

	variables := map[string]interface{}{
		"items":  githubv4.Int(c.PageItems),
		"query":  githubv4.String(query),
		"cursor": (*githubv4.String)(nil),
	}

	var q struct {
		Search search `graphql:"search(first: $items, after: $cursor type: REPOSITORY, query: $query)"`
	}
	var repos []map[string]string

	for {
		err := clnt.Query(ctx, &q, variables)
		if err != nil {
			return nil, err
		}
		for _, repo := range q.Search.Repos {
			m := make(map[string]string)
			m["name"] = repo.Repository.Name
			m["id"] = repo.Repository.ID
			m["url"] = repo.Repository.URL
			m["branch"] = repo.Repository.DefaultBranchRef.Name
			repos = append(repos, m)
		}
		if !q.Search.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.String(q.Search.PageInfo.EndCursor)
	}
	return repos, nil
}
