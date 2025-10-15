package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the loaded .scribe.yml file.
type Config struct {
    Sections     []Section `yaml:"sections" mapstructure:"sections"`
    IgnoreScopes []string  `yaml:"ignore_scopes" mapstructure:"ignore_scopes"`
}

// Section defines a single category in the final changelog.
type Section struct {
    Title string   `yaml:"title" mapstructure:"title"`
    Types []string `yaml:"types" mapstructure:"types"`
}

// Default returns the default configuration when no .scribe.yml is present.
func Default() *Config {
    return &Config{
        Sections: []Section{
            {Title: "Breaking Changes", Types: []string{}},
            {Title: "New Features", Types: []string{"feat"}},
            {Title: "Bug Fixes", Types: []string{"fix"}},
        },
        IgnoreScopes: []string{},
    }
}

// Load reads the configuration from <repoPath>/.scribe.yml. If the file does not
// exist, it returns the default configuration.
func Load(repoPath string) (*Config, error) {
    if repoPath == "" {
        return nil, errors.New("repoPath is required")
    }

    def := Default()

    // Look for either .scribe.yml or .scribe.yaml explicitly
    candidates := []string{
        filepath.Join(repoPath, ".scribe.yml"),
        filepath.Join(repoPath, ".scribe.yaml"),
    }

    var configFile string
    for _, c := range candidates {
        if info, err := os.Stat(c); err == nil && !info.IsDir() {
            configFile = c
            break
        }
    }

    if configFile == "" {
        return def, nil
    }

    v := viper.New()
    v.SetConfigFile(configFile)

    if err := v.ReadInConfig(); err != nil {
        return nil, err
    }
    cfg := &Config{}
    if err := v.Unmarshal(cfg); err != nil {
        return nil, err
    }
    // Merge defaults for any missing fields
    if len(cfg.Sections) == 0 {
        cfg.Sections = def.Sections
    }
    if cfg.IgnoreScopes == nil {
        cfg.IgnoreScopes = def.IgnoreScopes
    }
    return cfg, nil
}


