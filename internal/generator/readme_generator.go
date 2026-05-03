package generator

import (
	"log"
	"os"
	"strings"
)

func generateHeader() string {
	return "# Hey!\nI'm currently a Full-Stack developer.\n\n"
}

func generateTools() string {
	tools := []string{
		"ts", "go", "nodejs", "html", "css",
		"react", "express", "jest", "prisma",
		"tailwind", "postgres",
	}
	return "## Tools\n[![My Tools](https://skillicons.dev/icons?i=" + strings.Join(tools, ",") + ")](https://skillicons.dev)\n\n"
}

func generateStats(svgDir string) string {
	entries, err := os.ReadDir(svgDir)
	if err != nil {
		log.Println("SVGs generated directory not yet generated")
		return ""
	}

	var stats strings.Builder
	for _, e := range entries {
		stats.WriteString("<img src='")
		stats.WriteString(svgDir)
		stats.WriteByte('/')
		stats.WriteString(e.Name())
		stats.WriteString("' width='100%' />\n")
	}

	return stats.String()
}

func generateFooter() string {
	return "\n## Portfolio\n<a href='https://hsn-portfolio.pages.dev/'>Link to portfolio</a>"
}

func WriteReadme(svgDir string) error {
	f, err := os.OpenFile("README.md", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	f.WriteString(generateHeader())
	f.WriteString(generateTools())
	f.WriteString(generateStats(svgDir))
	f.WriteString(generateFooter())
	return nil
}
