package stats

import (
	"sync"

	"github.com/sp41414/sp41414/internal/client"
)

type StatsAggregator struct {
	Client              *client.Client
	LanguageClient      *client.LanguageClient
	ContributionsClient *client.ContributionsClient
}

func NewStatsAggregator(c *client.Client) *StatsAggregator {
	return &StatsAggregator{
		Client:              c,
		LanguageClient:      client.NewLanguageClient(c),
		ContributionsClient: client.NewContributionsClient(c),
	}
}

func (s *StatsAggregator) FetchStats() (*GitHubStats, error) {
	repos, err := s.Client.FetchRepos()
	if err != nil {
		return nil, err
	}

	var (
		wg                sync.WaitGroup
		languageStats     []client.LanguageStats
		contributionStats client.ContributionsStats
	)

	wg.Go(func() {
		languageStats = s.LanguageClient.CalculateLanguageStats(repos)
	})
	wg.Go(func() {
		contributionStats = s.ContributionsClient.CalculateContributionsStats()
	})

	wg.Wait()

	return &GitHubStats{
		Username:      s.Client.Username,
		Languages:     languageStats,
		Contributions: contributionStats,
		TotalRepos:    len(repos),
	}, nil
}
