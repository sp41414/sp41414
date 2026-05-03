package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sp41414/sp41414/internal/config"
	"github.com/sp41414/sp41414/internal/generator"
	"github.com/sp41414/sp41414/internal/stats"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %w", err)
	}

	githubToken := os.Getenv("GITHUB_TOKEN")
	if !strings.HasPrefix(strings.TrimSpace(githubToken), "ghp_") {
		log.Fatalf("GITHUB_TOKEN is not a valid github token, Example token: ghp_xxxxxxxxxxxx")
	}
	profileName := os.Getenv("GITHUB_USERNAME")
	generatedDir := "generated"
	svgWidth := 495

	var theme config.Theme
	if os.Getenv("THEME") == "dark" {
		theme = config.Theme{
			Bg:            "#1c1410",
			Border:        "#3d2f24",
			Text:          "#f0e6d3",
			TextSecondary: "#8c7560",
			Accent:        "#e0956e",
			Success:       "#7aad6e",
			Warning:       "#e0b84a",
			Err:           "#c95f4e",
		}
	} else {
		theme = config.Theme{
			Bg:            "#faecc9",
			Border:        "#ddc6a8",
			Text:          "#3d2a1a",
			TextSecondary: "#8c7560",
			Accent:        "#c46e25",
			Success:       "#5a8f5a",
			Warning:       "#c99a2e",
			Err:           "#b84c3a",
		}
	}

	config := config.NewConfig(githubToken, profileName, generatedDir, svgWidth, theme)
	stats := stats.NewStats(config)
	fmt.Println("Fetching github stats...")

	err = stats.FetchStats()
	if err != nil {
		log.Fatalf("Error fetching github stats: %w", err)
	}

	fmt.Println("Generating svgs...")
	err = stats.WriteStats()
	if err != nil {
		log.Fatalf("Error generating svgs: %w", err)
	}

	fmt.Println("Generating README...")
	err = generator.WriteReadme(generatedDir)
	if err != nil {
		log.Fatalf("Error generating README: %w", err)
	}

	fmt.Println("Done!")
}
