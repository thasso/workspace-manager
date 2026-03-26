// Package gitignore manages a marker-delimited section in .gitignore files.
package gitignore

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	BeginMarker = "# --- managed repos (do not edit below this line) ---"
	EndMarker   = "# --- end managed repos ---"
)

// Update writes the managed section of a .gitignore file, preserving any
// user content above the managed section.
func Update(path string, repoNames []string) error {
	var userContent string

	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading .gitignore: %w", err)
	}

	if err == nil {
		content := string(data)
		if idx := strings.Index(content, BeginMarker); idx >= 0 {
			userContent = content[:idx]
		} else {
			userContent = content
			if !strings.HasSuffix(userContent, "\n") {
				userContent += "\n"
			}
			userContent += "\n"
		}
	}

	var b strings.Builder
	b.WriteString(userContent)
	b.WriteString(BeginMarker + "\n")
	for _, name := range repoNames {
		fmt.Fprintf(&b, "%s/\n", name)
	}
	b.WriteString(EndMarker + "\n")

	// Atomic write
	dir := filepath.Dir(path)
	if dir == "" {
		dir = "."
	}
	tmp, err := os.CreateTemp(dir, ".gitignore-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmp.Name()

	if _, err := tmp.WriteString(b.String()); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("writing temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("closing temp file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("renaming temp file: %w", err)
	}
	return nil
}
