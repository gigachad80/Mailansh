package git

import (
	"os/exec"
	"regexp"
	"strings"
)

// GetNameForEmail retrieves the name associated with an email from git log
func GetNameForEmail(repoDir, email string) string {
	cmd := exec.Command("git", "log", "--format=%an", "--author="+email, "-1")
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// IsValidEmail validates if an email address is valid
func IsValidEmail(email string) bool {
	if len(email) < 5 || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}

	invalidPatterns := []string{
		"@example.com",
		"@test.com",
		"@localhost",
	}

	emailLower := strings.ToLower(email)
	for _, pattern := range invalidPatterns {
		if strings.Contains(emailLower, pattern) {
			return false
		}
	}

	return true
}

// IsNoReplyUsername checks if a username appears to be a noreply username
func IsNoReplyUsername(name string) bool {
	lowerName := strings.ToLower(name)
	return strings.Contains(lowerName, "noreply") ||
		regexp.MustCompile(`^[0-9]+$`).MatchString(name) ||
		regexp.MustCompile(`^[a-f0-9]{7,40}$`).MatchString(name)
}
