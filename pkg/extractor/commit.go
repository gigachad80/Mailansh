
package extractor

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"mailansh/pkg/color"
	"mailansh/pkg/contributor"
	"mailansh/pkg/git"
)

// ExtractFromCommitLog extracts contributors from git commit history
func ExtractFromCommitLog(ctx context.Context, repoDir string, contributorManager *contributor.Manager, quiet bool) {
	if !quiet {
		fmt.Printf("%s[INFO] Extracting from commit history...%s\n", color.Yellow, color.Reset)
	}

	cmd := exec.CommandContext(ctx, "git", "log", "--format=%an|%ae%n%cn|%ce", "--all")
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		if !quiet {
			fmt.Printf("%s[WARNING] Error fetching commit info: %v%s\n", color.Yellow, err, color.Reset)
		}
		return
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		email := strings.TrimSpace(parts[1])
		if git.IsValidEmail(email) && name != "" {
			contributorManager.Add(contributor.Contributor{Name: name, Email: email})
		}
	}
}
