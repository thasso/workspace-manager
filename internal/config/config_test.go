package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAndSaveRoundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "workspace.json")

	ws := &Workspace{
		Org:     "myorg",
		JiraURL: "https://jira.example.com",
		Repos: []Repo{
			{Name: "foo", RepoName: "foo-repo", Branch: "main", JiraKey: "FOO", Description: "Foo project"},
		},
		Projects: []Project{
			{Name: "Bar", JiraKey: "BAR"},
		},
	}

	if err := ws.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Org != ws.Org {
		t.Errorf("Org = %q, want %q", loaded.Org, ws.Org)
	}
	if loaded.JiraURL != ws.JiraURL {
		t.Errorf("JiraURL = %q, want %q", loaded.JiraURL, ws.JiraURL)
	}
	if len(loaded.Repos) != 1 {
		t.Fatalf("Repos count = %d, want 1", len(loaded.Repos))
	}
	if loaded.Repos[0].Name != "foo" {
		t.Errorf("Repos[0].Name = %q, want %q", loaded.Repos[0].Name, "foo")
	}
	if loaded.Repos[0].RepoName != "foo-repo" {
		t.Errorf("Repos[0].RepoName = %q, want %q", loaded.Repos[0].RepoName, "foo-repo")
	}
	if len(loaded.Projects) != 1 {
		t.Fatalf("Projects count = %d, want 1", len(loaded.Projects))
	}
}

func TestAddRepoDuplicate(t *testing.T) {
	ws := &Workspace{
		Org:   "myorg",
		Repos: []Repo{{Name: "foo", RepoName: "foo", Branch: "main"}},
	}

	err := ws.AddRepo(Repo{Name: "foo", RepoName: "foo2", Branch: "main"})
	if err == nil {
		t.Fatal("expected error for duplicate repo name")
	}
}

func TestAddRepoSuccess(t *testing.T) {
	ws := &Workspace{Org: "myorg", Repos: []Repo{}}

	if err := ws.AddRepo(Repo{Name: "bar", RepoName: "bar", Branch: "main"}); err != nil {
		t.Fatalf("AddRepo: %v", err)
	}
	if len(ws.Repos) != 1 {
		t.Fatalf("Repos count = %d, want 1", len(ws.Repos))
	}
}

func TestCloneURL(t *testing.T) {
	r := &Repo{Name: "foo", RepoName: "foo-repo"}
	url := r.CloneURL("myorg")
	want := "git@github.com:myorg/foo-repo.git"
	if url != want {
		t.Errorf("CloneURL = %q, want %q", url, want)
	}
}

func TestRepoNames(t *testing.T) {
	ws := &Workspace{
		Repos: []Repo{
			{Name: "a"},
			{Name: "b"},
		},
	}
	names := ws.RepoNames()
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Errorf("RepoNames = %v, want [a b]", names)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/workspace.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "workspace.json")
	os.WriteFile(path, []byte("not json"), 0644)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
