// Package config handles loading, saving, and manipulating workspace.json manifests.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Workspace is the top-level structure of a workspace.json manifest.
type Workspace struct {
	Org      string    `json:"org"`
	JiraURL  string    `json:"jira_url,omitempty"`
	Repos    []Repo    `json:"repos"`
	Projects []Project `json:"projects,omitempty"`
}

// Repo represents a git repository managed by the workspace.
type Repo struct {
	Name        string `json:"name"`
	RepoName    string `json:"repo"`
	Branch      string `json:"branch"`
	JiraKey     string `json:"jira_key,omitempty"`
	Description string `json:"description,omitempty"`
}

// Project represents a Jira-only project with no associated repository.
type Project struct {
	Name        string `json:"name"`
	JiraKey     string `json:"jira_key"`
	Description string `json:"description,omitempty"`
}

// CloneURL returns the SSH clone URL for a repo given the org.
func (r *Repo) CloneURL(org string) string {
	return fmt.Sprintf("git@github.com:%s/%s.git", org, r.RepoName)
}

// Load reads and parses a workspace.json file.
func Load(path string) (*Workspace, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading workspace config: %w", err)
	}
	var ws Workspace
	if err := json.Unmarshal(data, &ws); err != nil {
		return nil, fmt.Errorf("parsing workspace config: %w", err)
	}
	return &ws, nil
}

// Save writes the workspace config to the given path atomically.
func (ws *Workspace) Save(path string) error {
	data, err := json.MarshalIndent(ws, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling workspace config: %w", err)
	}
	data = append(data, '\n')

	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".workspace-*.json")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("writing temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("closing temp file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("renaming temp file: %w", err)
	}
	return nil
}

// AddRepo appends a repo to the workspace, returning an error if a repo with
// the same name already exists.
func (ws *Workspace) AddRepo(r Repo) error {
	for _, existing := range ws.Repos {
		if existing.Name == r.Name {
			return fmt.Errorf("repo %q already exists in workspace", r.Name)
		}
	}
	ws.Repos = append(ws.Repos, r)
	return nil
}

// RepoNames returns a slice of all repo directory names.
func (ws *Workspace) RepoNames() []string {
	names := make([]string, len(ws.Repos))
	for i, r := range ws.Repos {
		names[i] = r.Name
	}
	return names
}
