package git

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"mailansh/pkg/color"
)

// CreateTempDir creates a temporary directory for git operations
func CreateTempDir() string {
	dir, err := os.MkdirTemp("", fmt.Sprintf("git_clone_temp_%s-*", time.Now().Format("20060102_150405")))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError creating temp dir: %v%s\n", color.Red, err, color.Reset)
		os.Exit(1)
	}
	return dir
}

// CloneRepo clones a git repository to the specified directory
func CloneRepo(repoURL, repoDir string, quiet bool) {
	if !quiet {
		fmt.Printf("%s[INFO] Cloning repository...%s\n", color.Yellow, color.Reset)
	}
	cmd := exec.Command("git", "clone", repoURL, repoDir)

	if quiet {
		// Suppress all output in quiet mode
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%sError cloning repository: %v%s\n", color.Red, err, color.Reset)
		os.Exit(1)
	}
}

