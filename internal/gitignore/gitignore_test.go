package gitignore

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpdateCreatesNewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")

	if err := Update(path, []string{"repo-a", "repo-b"}); err != nil {
		t.Fatalf("Update: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)

	if !strings.Contains(content, BeginMarker) {
		t.Error("missing begin marker")
	}
	if !strings.Contains(content, EndMarker) {
		t.Error("missing end marker")
	}
	if !strings.Contains(content, "repo-a/") {
		t.Error("missing repo-a/")
	}
	if !strings.Contains(content, "repo-b/") {
		t.Error("missing repo-b/")
	}
}

func TestUpdatePreservesUserContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")

	userContent := "# My custom ignores\n*.log\n.env\n"
	os.WriteFile(path, []byte(userContent), 0644)

	if err := Update(path, []string{"repo-a"}); err != nil {
		t.Fatalf("Update: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)

	if !strings.HasPrefix(content, userContent) {
		t.Errorf("user content not preserved, got:\n%s", content)
	}
	if !strings.Contains(content, "repo-a/") {
		t.Error("missing repo-a/")
	}
}

func TestUpdateReplacesExistingManagedSection(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")

	initial := "# User stuff\n*.log\n" + BeginMarker + "\nold-repo/\n" + EndMarker + "\n"
	os.WriteFile(path, []byte(initial), 0644)

	if err := Update(path, []string{"new-repo"}); err != nil {
		t.Fatalf("Update: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)

	if strings.Contains(content, "old-repo/") {
		t.Error("old-repo should have been removed")
	}
	if !strings.Contains(content, "new-repo/") {
		t.Error("missing new-repo/")
	}
	if !strings.Contains(content, "# User stuff") {
		t.Error("user content not preserved")
	}
}

func TestUpdateIdempotent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")

	repos := []string{"a", "b"}

	Update(path, repos)
	data1, _ := os.ReadFile(path)

	Update(path, repos)
	data2, _ := os.ReadFile(path)

	if string(data1) != string(data2) {
		t.Error("Update is not idempotent")
	}
}
