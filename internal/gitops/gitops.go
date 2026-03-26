// Package gitops provides git operations via os/exec.
package gitops

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// RepoStatus holds the status information for a single repository.
type RepoStatus struct {
	Name          string
	CurrentBranch string
	DefaultBranch string
	IsCloned      bool
	IsDirty       bool
	Ahead         int
	Behind        int
	HasTracking   bool
}

// Clone clones a git repository to the given directory.
func Clone(url, branch, dir string) error {
	cmd := exec.Command("git", "clone", "--branch", branch, url, dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone: %s\n%s", err, string(out))
	}
	return nil
}

// InitSubmodules initializes and updates submodules recursively.
func InitSubmodules(dir string) error {
	cmd := exec.Command("git", "-C", dir, "submodule", "update", "--init", "--recursive")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git submodule update: %s\n%s", err, string(out))
	}
	return nil
}

// Pull runs git pull in the given directory. It returns true if changes were
// fetched, false if already up to date.
func Pull(dir string) (updated bool, err error) {
	cmd := exec.Command("git", "-C", dir, "pull")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("git pull: %s\n%s", err, string(out))
	}
	return !strings.Contains(string(out), "Already up to date"), nil
}

// UpdateSubmodules updates submodules after a pull.
func UpdateSubmodules(dir string) error {
	return InitSubmodules(dir)
}

// Status gathers the status of a repository.
func Status(dir, defaultBranch string) (*RepoStatus, error) {
	s := &RepoStatus{
		Name:          dir,
		DefaultBranch: defaultBranch,
		IsCloned:      true,
	}

	// Current branch
	branch, err := gitOutput(dir, "branch", "--show-current")
	if err != nil {
		return nil, err
	}
	s.CurrentBranch = branch
	if s.CurrentBranch == "" {
		s.CurrentBranch = "(detached)"
	}

	// Dirty state
	porcelain, err := gitOutput(dir, "status", "--porcelain")
	if err != nil {
		return nil, err
	}
	s.IsDirty = porcelain != ""

	// Ahead/behind
	tracking, err := gitOutput(dir, "rev-parse", "--abbrev-ref", "@{upstream}")
	if err == nil && tracking != "" {
		s.HasTracking = true

		ahead, err := gitOutput(dir, "rev-list", "--count", "@{upstream}..HEAD")
		if err == nil {
			s.Ahead, _ = strconv.Atoi(ahead)
		}
		behind, err := gitOutput(dir, "rev-list", "--count", "HEAD..@{upstream}")
		if err == nil {
			s.Behind, _ = strconv.Atoi(behind)
		}
	}

	return s, nil
}

func gitOutput(dir string, args ...string) (string, error) {
	fullArgs := append([]string{"-C", dir}, args...)
	cmd := exec.Command("git", fullArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %s\n%s", strings.Join(args, " "), err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}
