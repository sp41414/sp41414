package client

import (
	"fmt"
	"log"
	"slices"

	"github.com/go-enry/go-enry/v2"
	"github.com/google/go-github/v85/github"
)

type LanguageClient struct {
	Client *Client
}

type LanguageStats struct {
	Name       string
	Color      string
	Percentage float64
	Size       int
}

func NewLanguageClient(client *Client) *LanguageClient {
	return &LanguageClient{
		Client: client,
	}
}

func (c *LanguageClient) CalculateLanguageStats(repos []*github.Repository) []LanguageStats {
	languageBytes := make(map[string]int)

	count := 0
	for _, repo := range repos {
		languages, err := c.Client.FetchRepoLanguages(repo)
		if err != nil {
			log.Printf("Could not fetch languages for repo: %s", repo.GetName())
			continue
		}

		for language, bytes := range languages {
			if _, ok := languageBytes[language]; !ok {
				languageBytes[language] = 0 + bytes
			} else {
				languageBytes[language] += bytes
			}
		}

		count++
		fmt.Printf("Processed %d/%d repos\n", count, len(repos))
	}

	totalBytes := 0
	for _, bytes := range languageBytes {
		totalBytes += bytes
	}

	allLanguages := []LanguageStats{}
	for language, bytes := range languageBytes {
		allLanguages = append(allLanguages, LanguageStats{
			Name:       language,
			Size:       bytes,
			Percentage: (float64(bytes) / float64(totalBytes)) * 100,
			Color:      enry.GetColor(language),
		})
	}

	slices.SortFunc(allLanguages, func(a, b LanguageStats) int {
		if b.Percentage > a.Percentage {
			return 1
		}
		return -1
	})

	return allLanguages
}
