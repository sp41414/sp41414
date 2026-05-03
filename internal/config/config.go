package config

type Theme struct {
	Bg            string
	Border        string
	Text          string
	TextSecondary string
	Accent        string
	Success       string
	Warning       string
	Err           string
}

type Config struct {
	Theme        Theme
	GithubToken  string
	ProfileName  string
	GeneratedDir string
	SvgWidth     int
}

func NewConfig(githubToken, profileName, generatedDir string, svgWidth int, theme Theme) *Config {
	if githubToken == "" {
		panic("github token is required")
	}
	return &Config{
		Theme:        theme,
		GithubToken:  githubToken,
		ProfileName:  profileName,
		GeneratedDir: generatedDir,
		SvgWidth:     svgWidth,
	}
}
