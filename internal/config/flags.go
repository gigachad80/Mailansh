
package config

import "flag"

// All configuration options
type Config struct {
	Help          bool
	GitHubNoreply bool
	Popular       bool
	CustomDomain  bool
	OutputFile    string
	Quiet         bool
}

// ParseFlags parses command line flags and returns a Config
func ParseFlags() *Config {
	cfg := &Config{}

	flag.BoolVar(&cfg.Help, "h", false, "Show help menu")
	flag.BoolVar(&cfg.GitHubNoreply, "g", false, "Show only GitHub noreply emails")
	flag.BoolVar(&cfg.Popular, "p", false, "Show only popular email domains")
	flag.BoolVar(&cfg.CustomDomain, "cd", false, "Show only custom domains")
	flag.StringVar(&cfg.OutputFile, "o", "", "Save output to file (.csv or .txt)")
	flag.BoolVar(&cfg.Quiet, "q", false, "Quiet mode: no info or progress output")

	flag.Parse()
	return cfg
}
