package client

import (
	"context"

	"github.com/jferrl/go-githubauth"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Client struct {
	context  context.Context
	Client   *githubv4.Client
	Username string
}

type Repository struct {
	Name      string
	Owner     string
	Languages map[string]int
}

func NewClient(token, username string) *Client {
	tokenSource := githubauth.NewPersonalAccessTokenSource(token)
	ctx := context.Background()
	httpClient := oauth2.NewClient(ctx, tokenSource)
	client := githubv4.NewClient(httpClient)

	return &Client{
		Client:   client,
		Username: username,
		context:  ctx,
	}
}

func (c *Client) FetchRepos() ([]Repository, error) {
	var query struct {
		User struct {
			Repositories struct {
				Nodes []struct {
					Name  string
					Owner struct {
						Login string
					}
					Languages struct {
						Edges []struct {
							Size int
							Node struct {
								Name string
							}
						}
					} `graphql:"languages(first: 10, orderBy: {field: SIZE, direction: DESC})"`
				}
				PageInfo struct {
					EndCursor   githubv4.String
					HasNextPage bool
				}
			} `graphql:"repositories(first: 100, after: $cursor, ownerAffiliations: [OWNER], privacy: PUBLIC, isFork: false, orderBy: {field: UPDATED_AT, direction: DESC})"`
		} `graphql:"user(login: $username)"`
	}

	variables := map[string]interface{}{
		"username": githubv4.String(c.Username),
		"cursor":   (*githubv4.String)(nil),
	}

	var allRepos []Repository
	for {
		err := c.Client.Query(c.context, &query, variables)
		if err != nil {
			return allRepos, err
		}

		for _, node := range query.User.Repositories.Nodes {
			languages := make(map[string]int)
			for _, edge := range node.Languages.Edges {
				languages[edge.Node.Name] = edge.Size
			}
			allRepos = append(allRepos, Repository{
				Name:      node.Name,
				Owner:     node.Owner.Login,
				Languages: languages,
			})
		}

		if !query.User.Repositories.PageInfo.HasNextPage {
			break
		}
		variables["cursor"] = githubv4.NewString(query.User.Repositories.PageInfo.EndCursor)
	}

	return allRepos, nil
}
