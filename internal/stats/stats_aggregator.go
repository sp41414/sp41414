package stats

import (
	"github.com/sp41414/sp41414/internal/client"
)

type StatsAggregator struct {
	Client         *client.Client
	LanguageClient *client.LanguageClient
}

func NewStatsAggregator(c *client.Client) *StatsAggregator {
	return &StatsAggregator{
		Client:         c,
		LanguageClient: client.NewLanguageClient(c),
	}
}

func (s *StatsAggregator) FetchStats() (*GitHubStats, error) {
	repos, err := s.Client.FetchRepos()
	if err != nil {
		return nil, err
	}

	// make this concurrent when adding more clients
	languageStats := s.LanguageClient.CalculateLanguageStats(repos)

	return &GitHubStats{
		Username:   s.Client.Username,
		Languages:  languageStats,
		TotalRepos: len(repos),
	}, nil
}
