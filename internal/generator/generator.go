package generator

import (
	"github.com/sp41414/sp41414/internal/config"
)

type SVGGenerator struct {
	theme  *config.Theme
	width  int
	height int
}

type Generator interface {
	Generate(svg *SVGGenerator) error
}

func NewSVGGenerator(theme *config.Theme, width, height int) *SVGGenerator {
	return &SVGGenerator{
		theme:  theme,
		width:  width,
		height: height,
	}
}
