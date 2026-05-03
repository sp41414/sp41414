package client

import (
	"slices"

	"github.com/go-enry/go-enry/v2"
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

func (c *LanguageClient) CalculateLanguageStats(repos []Repository) []LanguageStats {
	languageBytes := make(map[string]int)

	for _, repo := range repos {
		for language, bytes := range repo.Languages {
			if _, ok := languageBytes[language]; !ok {
				languageBytes[language] = 0 + bytes
			} else {
				languageBytes[language] += bytes
			}
		}
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
