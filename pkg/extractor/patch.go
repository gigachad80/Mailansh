
package extractor

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"mailansh/pkg/color"
	"mailansh/pkg/contributor"
	"mailansh/pkg/git"
	"mailansh/pkg/platform"
	"mailansh/pkg/progress"
)

// ExtractFromPatchesConcurrently extracts contributors from patch content using worker pools
func ExtractFromPatchesConcurrently(ctx context.Context, repoDir string, platform platform.Platform, contributorManager *contributor.Manager, quiet bool) {
	if !quiet {
		fmt.Printf("%s[INFO] Extracting from patch content...%s\n", color.Yellow, color.Reset)
	}

	// Get all commit hashes
	cmd := exec.CommandContext(ctx, "git", "log", "--reverse", "--format=%H")
	cmd.Dir = repoDir
	output, err := cmd.Output()
	if err != nil {
		if !quiet {
			fmt.Printf("%s[WARNING] Error fetching commit hashes: %v%s\n", color.Yellow, err, color.Reset)
		}
		return
	}

	hashes := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(hashes) == 0 || hashes[0] == "" {
		return
	}

	// Worker pool configuration
	const numWorkers = 8

	// Create progress reporter
	progressReporter := progress.NewReporter(len(hashes), quiet)
	defer progressReporter.Stop()

	// Create channels for worker pool
	hashChan := make(chan string, numWorkers*2)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			processCommitWorker(ctx, workerID, repoDir, platform, hashChan, contributorManager, progressReporter)
		}(i)
	}

	// Send work to workers
	go func() {
		defer close(hashChan)
		for _, hash := range hashes {
			hash = strings.TrimSpace(hash)
			if hash == "" {
				continue
			}

			select {
			case hashChan <- hash:
			case <-ctx.Done():
				return
			}
		}
	}()

	wg.Wait()
}

// processCommitWorker processes individual commits in worker goroutines
func processCommitWorker(ctx context.Context, workerID int, repoDir string, platform platform.Platform, hashChan <-chan string, contributorManager *contributor.Manager, progressReporter *progress.Reporter) {
	emailRegex := platform.BuildEmailRegex()

	for hash := range hashChan {
		select {
		case <-ctx.Done():
			return
		default:
			// Get the patch content
			cmd := exec.CommandContext(ctx, "git", "show", hash)
			cmd.Dir = repoDir
			content, err := cmd.Output()
			if err != nil {
				progressReporter.Add(1)
				continue
			}

			// Extract emails from the patch content
			matches := emailRegex.FindAllString(string(content), -1)
			for _, email := range matches {
				email = strings.TrimSpace(email)
				if git.IsValidEmail(email) {
					name := git.GetNameForEmail(repoDir, email)
					if name == "" {
						name = "Unknown"
					}
					contributorManager.Add(contributor.Contributor{Name: name, Email: email})
				}
			}

			progressReporter.Add(1)
		}
	}
}
