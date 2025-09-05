
package extractor

import (
	"context"
	"sync"

	"mailansh/pkg/contributor"
	"mailansh/pkg/platform"
)

// Extractor handles contributor extraction from git repositories
type Extractor struct {
	platform platform.Platform
	quiet    bool
}

// NewExtractor creates a new extractor
func NewExtractor(platform platform.Platform, quiet bool) *Extractor {
	return &Extractor{
		platform: platform,
		quiet:    quiet,
	}
}

// ExtractConcurrently extracts contributors using concurrent processing
func (e *Extractor) ExtractConcurrently(ctx context.Context, repoDir string) []contributor.Contributor {
	contributorManager := contributor.NewManager()
	var wg sync.WaitGroup

	// Parallel extraction from different sources
	wg.Add(2)

	// Goroutine 1: Extract from commit log
	go func() {
		defer wg.Done()
		ExtractFromCommitLog(ctx, repoDir, contributorManager, e.quiet)
	}()

	// Goroutine 2: Extract from patchaes with worker pool
	go func() {
		defer wg.Done()
		ExtractFromPatchesConcurrently(ctx, repoDir, e.platform, contributorManager, e.quiet)
	}()

	wg.Wait()
	return contributorManager.GetAll()
}
