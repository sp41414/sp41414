package client

import (
	"log"
	"time"

	"github.com/shurcooL/githubv4"
)

type ContributionsClient struct {
	Client *Client
}

type ContributionsStats struct {
	TotalCommits           int
	TotalIssues            int
	TotalPRs               int
	TotalRepoContributions int
}

func NewContributionsClient(client *Client) *ContributionsClient {
	return &ContributionsClient{
		Client: client,
	}
}

func (c *ContributionsClient) CalculateContributionsStats() ContributionsStats {
	today := time.Now()
	oneYearAgo := today.AddDate(-1, 0, 0)

	var query struct {
		User struct {
			ContributionsCollection struct {
				TotalCommitContributions      int
				TotalIssueContributions       int
				TotalPullRequestContributions int
				TotalRepositoryContributions  int
				RestrictedContributionsCount  int
			} `graphql:"contributionsCollection(from: $from, to: $to)"`
		} `graphql:"user(login: $username)"`
	}

	variables := map[string]interface{}{
		"username": githubv4.String(c.Client.Username),
		"from":     githubv4.DateTime{Time: oneYearAgo},
		"to":       githubv4.DateTime{Time: today},
	}

	err := c.Client.Client.Query(c.Client.context, &query, variables)
	if err != nil {
		log.Printf("Error fetching contributions: %s", err)
		return ContributionsStats{}
	}

	col := query.User.ContributionsCollection

	return ContributionsStats{
		TotalCommits:           col.TotalCommitContributions + col.RestrictedContributionsCount,
		TotalIssues:            col.TotalIssueContributions,
		TotalPRs:               col.TotalPullRequestContributions,
		TotalRepoContributions: col.TotalRepositoryContributions,
	}
}
