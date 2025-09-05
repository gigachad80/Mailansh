
package output

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"mailansh/pkg/color"
	"mailansh/pkg/contributor"
)

// Formatter handles output formatting and file operations
type Formatter struct {
	quiet bool
}

// NewFormatter creates a new output formatter
func NewFormatter(quiet bool) *Formatter {
	return &Formatter{quiet: quiet}
}

// Display shows contributors in a formatted table
func (f *Formatter) Display(contributors []contributor.Contributor) {
	if len(contributors) == 0 {
		fmt.Printf("%s[INFO] No contributors found matching the specified criteria.%s\n", color.Yellow, color.Reset)
		return
	}

	if f.quiet {
		// In quiet mode, just output name,email pairs
		for _, c := range contributors {
			fmt.Printf("%s,%s\n", c.Name, c.Email)
		}
	} else {
		fmt.Printf("%s%-30s %s%s%s\n", color.Cyan, "NAME", color.Magenta, "EMAIL", color.Reset)
		fmt.Printf("%s%s%s\n", color.Blue, strings.Repeat("=", 70), color.Reset)
		for _, contributor := range contributors {
			fmt.Printf("%s%-30s %s%s%s\n", color.Green, contributor.Name, color.Yellow, contributor.Email, color.Reset)
		}
	}
}

// SaveToFile saves contributors to a CSV or TXT file
func (f *Formatter) SaveToFile(contributors []contributor.Contributor, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	if strings.HasSuffix(strings.ToLower(filename), ".csv") {
		// Save as CSV format
		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Write header
		if err := writer.Write([]string{"Name", "Email"}); err != nil {
			return fmt.Errorf("failed to write CSV header: %v", err)
		}

		// Write contributor data
		for _, c := range contributors {
			if err := writer.Write([]string{c.Name, c.Email}); err != nil {
				return fmt.Errorf("failed to write CSV row: %v", err)
			}
		}
	} else {
		// Save as TXT format (comma-separated)
		for _, c := range contributors {
			if _, err := file.WriteString(fmt.Sprintf("%s,%s\n", c.Name, c.Email)); err != nil {
				return fmt.Errorf("failed to write to file: %v", err)
			}
		}
	}

	return nil
}
