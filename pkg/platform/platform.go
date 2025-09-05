
package platform

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Platform represents different git hosting platforms
type Platform int

const (
	GitHub Platform = iota
	GitLab
	Gitea
	Bitbucket
)

// String returns the string representation of the platform
func (p Platform) String() string {
	switch p {
	case GitHub:
		return "GitHub"
	case GitLab:
		return "GitLab"
	case Gitea:
		return "Gitea"
	case Bitbucket:
		return "Bitbucket"
	default:
		return "Unknown"
	}
}

// ParseRepoURL parses a repository URL and returns the platform and repo name
func ParseRepoURL(repoURL string) (Platform, string, error) {
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return 0, "", fmt.Errorf("invalid URL format")
	}

	host := strings.ToLower(parsedURL.Host)
	var platform Platform

	switch {
	case strings.Contains(host, "github.com"):
		platform = GitHub
	case strings.Contains(host, "gitlab.com") || strings.Contains(host, "gitlab."):
		platform = GitLab
	case strings.Contains(host, "gitea.com") || strings.Contains(host, "gitea."):
		platform = Gitea
	case strings.Contains(host, "bitbucket.org") || strings.Contains(host, "bitbucket."):
		platform = Bitbucket
	default:
		if strings.Contains(repoURL, "/tree/") || strings.Contains(repoURL, "/blob/") {
			platform = GitLab
		} else if strings.Contains(repoURL, "/src/") {
			platform = Gitea
		} else {
			return 0, "", fmt.Errorf("unsupported platform: %s", host)
		}
	}

	repoName, err := extractRepoName(parsedURL, platform)
	if err != nil {
		return 0, "", err
	}

	return platform, repoName, nil
}

// BuildEmailRegex builds a regex pattern for extracting emails based on platform
func (p Platform) BuildEmailRegex() *regexp.Regexp {
	patterns := []string{
		`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`,
	}

	switch p {
	case GitHub:
		patterns = append(patterns,
			`[0-9]+\+[a-zA-Z0-9._%+\-]+@users\.noreply\.github\.com`,
			`[a-zA-Z0-9._%+\-]+@users\.noreply\.github\.com`,
		)
	case GitLab:
		patterns = append(patterns,
			`[0-9]+\-[a-zA-Z0-9._%+\-]+@users\.noreply\.gitlab\.com`,
			`[a-zA-Z0-9._%+\-]+@users\.noreply\.gitlab\.com`,
		)
	case Gitea:
		patterns = append(patterns,
			`[a-zA-Z0-9._%+\-]+@noreply\.gitea\.[a-zA-Z0-9.\-]+`,
			`[a-zA-Z0-9._%+\-]+@users\.noreply\.gitea\.com`,
		)
	case Bitbucket:
		patterns = append(patterns,
			`[a-zA-Z0-9._%+\-]+@noreply\.bitbucket\.org`,
		)
	}

	combinedPattern := fmt.Sprintf(`\b(%s)\b`, strings.Join(patterns, "|"))
	return regexp.MustCompile(combinedPattern)
}

// extractRepoName extracts repository name from parsed URL
func extractRepoName(parsedURL *url.URL, platform Platform) (string, error) {
	path := strings.Trim(parsedURL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("cannot parse repository name from URL")
	}

	switch platform {
	case GitHub, Gitea:
		return strings.TrimSuffix(parts[1], ".git"), nil
	case GitLab:
		if len(parts) >= 2 {
			return strings.TrimSuffix(parts[len(parts)-1], ".git"), nil
		}
		return "", fmt.Errorf("cannot parse GitLab repository name")
	case Bitbucket:
		return strings.TrimSuffix(parts[1], ".git"), nil
	default:
		return "", fmt.Errorf("unsupported platform")
	}
}
