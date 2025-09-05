
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mailansh/internal/config"
	"mailansh/pkg/color"
	"mailansh/pkg/contributor"
	"mailansh/pkg/extractor"
	"mailansh/pkg/git"
	"mailansh/pkg/output"
	"mailansh/pkg/platform"
)

func main() {
	cfg := config.ParseFlags()

	if cfg.Help || flag.NArg() < 1 {
		showHelp()
		return
	}

	repoURL := strings.TrimSpace(flag.Arg(0))
	platform, repoName, err := platform.ParseRepoURL(repoURL)
	if err != nil {
		if !cfg.Quiet {
			fmt.Fprintf(os.Stderr, "%sInvalid repository URL: %v%s\n", color.Red, err, color.Reset)
		}
		os.Exit(1)
	}

	if !cfg.Quiet {
		fmt.Printf("%s[INFO] Detected platform: %s%s\n", color.Blue, platform, color.Reset)
		fmt.Printf("%s[INFO] Repository: %s%s\n", color.Blue, repoName, color.Reset)
	}

	// Setup directories
	tempCloneDir := git.CreateTempDir()
	defer os.RemoveAll(tempCloneDir)
	repoDir := filepath.Join(tempCloneDir, repoName)

	// Clone the repository
	git.CloneRepo(repoURL, repoDir, cfg.Quiet)

	// Extract contributors with concurrency
	contributorExtractor := extractor.NewExtractor(platform, cfg.Quiet)
	contributors := contributorExtractor.ExtractConcurrently(context.Background(), repoDir)

	// Filter contributors based on flags
	filteredContributors := contributor.FilterContributors(contributors, cfg)

	if !cfg.Quiet {
		fmt.Printf("%s[INFO] Found %d unique contributors (after filtering):%s\n\n", color.Green, len(filteredContributors), color.Reset)
	}

	// Sort contributors for consistent output
	contributor.SortByName(filteredContributors)

	// Save to file or display on console
	formatter := output.NewFormatter(cfg.Quiet)
	if cfg.OutputFile != "" {
		err := formatter.SaveToFile(filteredContributors, cfg.OutputFile)
		if err != nil {
			if !cfg.Quiet {
				fmt.Fprintf(os.Stderr, "%sError saving to file: %v%s\n", color.Red, err, color.Reset)
			}
			os.Exit(1)
		}
		if !cfg.Quiet {
			fmt.Printf("%s[SUCCESS] Output saved to %s%s\n", color.Green, cfg.OutputFile, color.Reset)
		}
	} else {
		formatter.Display(filteredContributors)
	}
}

func showHelp() {
	fmt.Printf("%s=== Git Contributors Extractor ===%s\n\n", color.Cyan, color.Reset)
	fmt.Printf("%sUsage:%s\n", color.Yellow, color.Reset)
	fmt.Printf("  %smailansh [flags] <repository-url>%s\n\n", color.Blue, color.Reset)
	fmt.Printf("%sFlags:%s\n", color.Yellow, color.Reset)
	fmt.Printf("  %s-h%s          Show this help menu\n", color.Green, color.Reset)
	fmt.Printf("  %s-g%s          Show only GitHub noreply emails (users.noreply.github.com)\n", color.Green, color.Reset)
	fmt.Printf("  %s-p%s          Show only popular email domains (gmail, protonmail, hotmail, outlook, yahoo)\n", color.Green, color.Reset)
	fmt.Printf("  %s-cd%s         Show only custom domains (excluding popular and noreply domains)\n", color.Green, color.Reset)
	fmt.Printf("  %s-o <file>%s   Save output to file (.csv or .txt format)\n", color.Green, color.Reset)
	fmt.Printf("  %s-q%s          Quiet mode: no progress or info output, only final output\n", color.Green, color.Reset)
	fmt.Printf("\n%sSupported platforms:%s GitHub, GitLab, Gitea, Bitbucket\n", color.Yellow, color.Reset)
	fmt.Printf("\n%sExamples:%s\n", color.Yellow, color.Reset)
	fmt.Printf("  mailansh https://github.com/user/repo\n")
	fmt.Printf("  mailansh -g https://github.com/user/repo\n")
	fmt.Printf("  mailansh -p -o contributors.csv https://github.com/user/repo\n")
	fmt.Printf("  mailansh -cd -q https://github.com/user/repo\n")
}
