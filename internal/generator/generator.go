package generator

import (
	"github.com/sp41414/sp41414/internal/config"
)

type SVGGenerator struct {
	theme *config.Theme
	width int
}

func NewSVGGenerator(theme *config.Theme, width int) *SVGGenerator {
	return &SVGGenerator{
		theme: theme,
		width: width,
	}
}
