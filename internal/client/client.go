package client

import (
	"context"

	"github.com/google/go-github/v85/github"
	"github.com/jferrl/go-githubauth"
	"golang.org/x/oauth2"
)

type Client struct {
	context  context.Context
	Client   *github.Client
	Username string
}

func NewClient(token, username string) *Client {
	tokenSource := githubauth.NewPersonalAccessTokenSource(token)
	ctx := context.Background()
	httpClient := oauth2.NewClient(ctx, tokenSource)
	client := github.NewClient(httpClient)

	return &Client{
		Client:   client,
		Username: username,
		context:  ctx,
	}
}

func (c *Client) FetchProfile() (*github.User, error) {
	user, _, err := c.Client.Users.Get(c.context, c.Username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Client) FetchRepos() ([]*github.Repository, error) {
	var allRepos []*github.Repository
	opts := &github.RepositoryListByAuthenticatedUserOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Affiliation: "owner",
		Sort:        "updated",
	}
	iter := c.Client.Repositories.ListByAuthenticatedUserIter(c.context, opts)
	for repo, err := range iter {
		if err != nil {
			return allRepos, err
		}
		allRepos = append(allRepos, repo)
	}

	return allRepos, nil
}

func (c *Client) FetchRepoLanguages(repo *github.Repository) (map[string]int, error) {
	languages, _, err := c.Client.Repositories.ListLanguages(c.context, repo.GetOwner().GetLogin(), repo.GetName())
	if err != nil {
		return nil, err
	}
	return languages, nil
}
