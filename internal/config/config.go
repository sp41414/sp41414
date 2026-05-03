package config

type Config struct {
	GithubToken  string
	ProfileName  string
	GeneratedDir string
}

func NewConfig(githubToken, profileName, generatedDir string) *Config {
	if githubToken == "" {
		panic("github token is required")
	}
	return &Config{
		GithubToken:  githubToken,
		ProfileName:  profileName,
		GeneratedDir: generatedDir,
	}
}
