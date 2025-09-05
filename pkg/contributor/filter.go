
package contributor

import (
	"regexp"
	"strings"

	"mailansh/internal/config"
)

// FilterContributors filters contributors based on configuration flags
func FilterContributors(contributors []Contributor, cfg *config.Config) []Contributor {
	// If no filter flags are set, return all contributors
	if !cfg.GitHubNoreply && !cfg.Popular && !cfg.CustomDomain {
		return contributors
	}

	var filtered []Contributor

	// Define popular email domains
	popularDomains := []string{
		"gmail.com", "protonmail.com", "protonmail.ch", "pm.me", "aol.com",
		"hotmail.com", "outlook.com", "live.com", "msn.com", "proton.me",
		"yahoo.com", "yahoo.co.uk", "ymail.com", "outlook.in",
		"aol.com", "icloud.com", "me.com", "mac.com",
	}

	// GitHub noreply regex pattern
	githubNoreplyRegex := regexp.MustCompile(`^([0-9]+\+)?[a-zA-Z0-9._%+\-]+@users\.noreply\.github\.com$`)

	for _, c := range contributors {
		emailLower := strings.ToLower(c.Email)

		// Check for GitHub noreply emails
		if cfg.GitHubNoreply && githubNoreplyRegex.MatchString(emailLower) {
			filtered = append(filtered, c)
			continue
		}

		// Check for popular domains
		if cfg.Popular {
			isPopular := false
			for _, domain := range popularDomains {
				if strings.HasSuffix(emailLower, "@"+domain) {
					isPopular = true
					break
				}
			}
			if isPopular {
				filtered = append(filtered, c)
				continue
			}
		}

		// Check for custom domains
		if cfg.CustomDomain {
			isPopular := false
			for _, domain := range popularDomains {
				if strings.HasSuffix(emailLower, "@"+domain) {
					isPopular = true
					break
				}
			}

			isNoreply := githubNoreplyRegex.MatchString(emailLower) ||
				strings.Contains(emailLower, "noreply") ||
				strings.Contains(emailLower, "no-reply")

			if !isPopular && !isNoreply && isValidCustomDomain(emailLower) {
				filtered = append(filtered, c)
			}
		}
	}

	return filtered
}

// isValidCustomDomain checks if an email has a valid custom domain
func isValidCustomDomain(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	domain := parts[1]
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$`)
	return domainRegex.MatchString(domain)
}
