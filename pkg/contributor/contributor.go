
package contributor

import (
	"sort"
	"strings"
	"sync"
)

// Contributor represents a git contributor
type Contributor struct {
	Name  string
	Email string
	mu    sync.Mutex
}

// SortByName sorts contributors alphabetically by name
func SortByName(contributors []Contributor) {
	sort.Slice(contributors, func(i, j int) bool {
		return strings.ToLower(contributors[i].Name) < strings.ToLower(contributors[j].Name)
	})
}
