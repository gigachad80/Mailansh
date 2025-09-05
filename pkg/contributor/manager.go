
package contributor

import (
	"regexp"
	"strings"
	"sync"
)

// Manager manages a collection of contributors with thread-safety
type Manager struct {
	contributors map[string]Contributor
	mu           sync.RWMutex
}

// NewManager creates a new contributor manager
func NewManager() *Manager {
	return &Manager{
		contributors: make(map[string]Contributor),
	}
}

// Add adds a contributor to the manager
func (cm *Manager) Add(contributor Contributor) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if existing, exists := cm.contributors[contributor.Email]; exists {
		// Prefer non-noreply names
		if isNoReplyUsername(existing.Name) && !isNoReplyUsername(contributor.Name) {
			cm.contributors[contributor.Email] = contributor
		}
	} else {
		cm.contributors[contributor.Email] = contributor
	}
}

// GetAll returns all contributors as a slice
func (cm *Manager) GetAll() []Contributor {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var result []Contributor
	for _, contributor := range cm.contributors {
		result = append(result, contributor)
	}
	return result
}

// isNoReplyUsername checks if a name looks like a noreply username
func isNoReplyUsername(name string) bool {
	lowerName := strings.ToLower(name)
	return strings.Contains(lowerName, "noreply") ||
		regexp.MustCompile(`^[0-9]+$`).MatchString(name) ||
		regexp.MustCompile(`^[a-f0-9]{7,40}$`).MatchString(name)
}
