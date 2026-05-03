package generator

import (
	"fmt"
	"os"
	"path/filepath"

	svg "github.com/ajstarks/svgo"
	"github.com/sp41414/sp41414/internal/client"
)

type ContributionsGenerator struct {
	generator *SVGGenerator
	stats     client.ContributionsStats
	outputDir string
}

func NewContributionsGenerator(gen *SVGGenerator, stats client.ContributionsStats, outputDir string) *ContributionsGenerator {
	return &ContributionsGenerator{
		generator: gen,
		stats:     stats,
		outputDir: outputDir,
	}
}

func (g *ContributionsGenerator) Generate() error {
	path := filepath.Join(g.outputDir, "contributions.svg")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create contributions svg: %w", err)
	}
	defer f.Close()

	const (
		padding     = 20
		titleHeight = 22
		cornerR     = 12
		rowHeight   = 32
	)

	type statRow struct {
		label string
		value int
	}
	rows := []statRow{
		{"Commits", g.stats.TotalCommits},
		{"Pull Requests", g.stats.TotalPRs},
		{"Issues", g.stats.TotalIssues},
		{"Repositories", g.stats.TotalRepoContributions},
	}

	theme := g.generator.theme
	width := g.generator.width
	height := padding + titleHeight + padding + (len(rows) * rowHeight) + padding

	canvas := svg.New(f)
	canvas.Startraw(fmt.Sprintf(
		`width="%d" height="%d" viewBox="0 0 %d %d" fill="none" xmlns="http://www.w3.org/2000/svg"`,
		width, height, width, height,
	))

	canvas.Roundrect(0, 0, width, height, cornerR, cornerR,
		fmt.Sprintf(`fill="%s" stroke="%s" stroke-width="1"`, theme.Bg, theme.Border))

	canvas.Text(padding, padding+14, "Contributions",
		fmt.Sprintf(`font-family="'Segoe UI',sans-serif" font-size="14" font-weight="600" fill="%s"`, theme.Text))

	rowsY := padding + titleHeight + padding
	for i, row := range rows {
		ry := rowsY + i*rowHeight
		if i > 0 {
			canvas.Line(padding, ry, width-padding, ry,
				fmt.Sprintf(`stroke="%s" stroke-width="1"`, theme.Border))
		}
		canvas.Text(padding, ry+20, row.label,
			fmt.Sprintf(`font-family="'Segoe UI',sans-serif" font-size="12" fill="%s"`, theme.TextSecondary))
		canvas.Text(width-padding, ry+20, fmt.Sprintf("%d", row.value),
			fmt.Sprintf(`font-family="'Segoe UI',sans-serif" font-size="12" font-weight="600" fill="%s" text-anchor="end"`, theme.Text))
	}

	canvas.End()
	return nil
}
