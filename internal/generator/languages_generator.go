package generator

import (
	"fmt"
	"os"
	"path/filepath"

	svg "github.com/ajstarks/svgo"
	"github.com/sp41414/sp41414/internal/client"
)

type LanguagesGenerator struct {
	generator *SVGGenerator
	stats     []client.LanguageStats
	outputDir string
}

func NewLanguagesGenerator(gen *SVGGenerator, stats []client.LanguageStats, outputDir string) *LanguagesGenerator {
	return &LanguagesGenerator{
		generator: gen,
		stats:     stats,
		outputDir: outputDir,
	}
}

func (g *LanguagesGenerator) Generate() error {
	path := filepath.Join(g.outputDir, "languages.svg")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create languages svg: %w", err)
	}
	defer f.Close()

	const (
		padding   = 20
		barHeight = 8
		rowHeight = 28
		dotSize   = 10
		cornerR   = 12
	)

	langs := g.stats
	if len(langs) > 8 {
		langs = langs[:8]
	}

	theme := g.generator.theme
	width := g.generator.width
	legendRows := (len(langs) + 1) / 2
	height := padding + 22 + padding/2 + barHeight + padding + (legendRows * rowHeight) + padding

	canvas := svg.New(f)
	canvas.Startraw(fmt.Sprintf(
		`width="%d" height="%d" viewBox="0 0 %d %d" fill="none" xmlns="http://www.w3.org/2000/svg"`,
		width, height, width, height,
	))

	canvas.Roundrect(0, 0, width, height, cornerR, cornerR,
		fmt.Sprintf(`fill="%s" stroke="%s" stroke-width="1"`, theme.Bg, theme.Border))

	canvas.Text(padding, padding+14, "Languages",
		fmt.Sprintf(`font-family="'Segoe UI',sans-serif" font-size="14" font-weight="600" fill="%s"`, theme.Text))

	barY := padding + 22 + padding/2
	barWidth := width - padding*2
	x := padding
	for _, lang := range langs {
		segW := int((lang.Percentage / 100) * float64(barWidth))
		if segW < 2 {
			segW = 2
		}
		canvas.Rect(x, barY, segW, barHeight,
			fmt.Sprintf(`fill="%s" rx="2"`, lang.Color))
		x += segW + 2
	}

	legendY := barY + barHeight + padding
	col0X := padding
	col1X := width/2 + padding/2

	for i, lang := range langs {
		col := i % 2
		row := i / 2
		lx := col0X
		if col == 1 {
			lx = col1X
		}
		ly := legendY + row*rowHeight

		canvas.Circle(lx+dotSize/2, ly+dotSize/2, dotSize/2,
			fmt.Sprintf(`fill="%s"`, lang.Color))

		canvas.Text(lx+dotSize+8, ly+dotSize/2+4, lang.Name,
			fmt.Sprintf(`font-family="'Segoe UI',sans-serif" font-size="12" fill="%s"`, theme.Text))

		pctX := col1X - padding/2
		if col == 1 {
			pctX = width - padding
		}
		canvas.Text(pctX, ly+dotSize/2+4,
			fmt.Sprintf("%.2f%%", lang.Percentage),
			fmt.Sprintf(`font-family="'Segoe UI',sans-serif" font-size="11" fill="%s" text-anchor="end"`, theme.TextSecondary))
	}

	canvas.End()
	return nil
}
