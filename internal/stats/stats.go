package stats

import (
	"fmt"
	"os"
	"sync"

	"github.com/sp41414/sp41414/internal/client"
	"github.com/sp41414/sp41414/internal/config"
	"github.com/sp41414/sp41414/internal/generator"
)

type GitHubStats struct {
	Languages     []client.LanguageStats
	Contributions client.ContributionsStats
	Username      string
	TotalRepos    int
}

type Stats struct {
	config *config.Config
	client *client.Client
	agg    *StatsAggregator
	stats  *GitHubStats
}

func NewStats(c *config.Config) *Stats {
	return &Stats{
		config: c,
	}
}

func (s *Stats) FetchStats() error {
	s.client = client.NewClient(s.config.GithubToken, s.config.ProfileName)
	s.agg = NewStatsAggregator(s.client)
	stats, err := s.agg.FetchStats()
	if err != nil {
		return err
	}
	s.stats = stats
	return nil
}

func (s *Stats) WriteStats() error {
	if s.stats == nil {
		return fmt.Errorf("stats not fetched yet")
	}

	if _, err := os.Stat(s.config.GeneratedDir); os.IsNotExist(err) {
		err := os.Mkdir(s.config.GeneratedDir, 0755)
		if err != nil {
			return err
		}
	}

	svgGenerator := generator.NewSVGGenerator(&s.config.Theme, s.config.SvgWidth)

	langGen := generator.NewLanguagesGenerator(svgGenerator, s.stats.Languages, s.config.GeneratedDir)
	contribGen := generator.NewContributionsGenerator(svgGenerator, s.stats.Contributions, s.config.GeneratedDir)

	var wg sync.WaitGroup
	errs := make(chan error, 2)
	wg.Go(func() {
		errs <- langGen.Generate()
	})
	wg.Go(func() {
		errs <- contribGen.Generate()
	})

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
