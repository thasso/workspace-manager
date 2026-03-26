package gitops

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const testBranch = "main"

// createTestRepo creates a bare git repo and a working clone for testing.
func createTestRepo(t *testing.T) (bareDir, workDir string) {
	t.Helper()
	base := t.TempDir()

	bareDir = filepath.Join(base, "bare.git")
	if err := exec.Command("git", "init", "--bare", "--initial-branch="+testBranch, bareDir).Run(); err != nil {
		t.Fatalf("git init --bare: %v", err)
	}

	workDir = filepath.Join(base, "work")
	if err := exec.Command("git", "clone", bareDir, workDir).Run(); err != nil {
		t.Fatalf("git clone: %v", err)
	}

	// Configure git user for commits
	exec.Command("git", "-C", workDir, "config", "user.email", "test@test.com").Run()
	exec.Command("git", "-C", workDir, "config", "user.name", "Test").Run()

	// Create an initial commit and push
	dummyFile := filepath.Join(workDir, "README.md")
	os.WriteFile(dummyFile, []byte("# Test\n"), 0644)
	exec.Command("git", "-C", workDir, "add", ".").Run()
	exec.Command("git", "-C", workDir, "commit", "-m", "initial commit").Run()
	exec.Command("git", "-C", workDir, "push", "origin", testBranch).Run()

	return bareDir, workDir
}

// defaultBranch returns the current branch name in the work dir.
func defaultBranch(t *testing.T, dir string) string {
	t.Helper()
	out, err := exec.Command("git", "-C", dir, "branch", "--show-current").Output()
	if err != nil {
		t.Fatalf("get branch: %v", err)
	}
	return strings.TrimSpace(string(out))
}

func TestStatusClean(t *testing.T) {
	_, workDir := createTestRepo(t)
	branch := defaultBranch(t, workDir)

	s, err := Status(workDir, branch)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}

	if s.IsDirty {
		t.Error("expected clean repo")
	}
	if s.CurrentBranch != branch {
		t.Errorf("CurrentBranch = %q, want %q", s.CurrentBranch, branch)
	}
	if !s.HasTracking {
		t.Error("expected tracking branch")
	}
	if s.Ahead != 0 {
		t.Errorf("Ahead = %d, want 0", s.Ahead)
	}
}

func TestStatusDirty(t *testing.T) {
	_, workDir := createTestRepo(t)
	branch := defaultBranch(t, workDir)

	os.WriteFile(filepath.Join(workDir, "dirty.txt"), []byte("dirty"), 0644)

	s, err := Status(workDir, branch)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}

	if !s.IsDirty {
		t.Error("expected dirty repo")
	}
}

func TestCloneAndStatus(t *testing.T) {
	bareDir, _ := createTestRepo(t)

	cloneDir := filepath.Join(t.TempDir(), "cloned")
	if err := Clone(bareDir, testBranch, cloneDir); err != nil {
		t.Fatalf("Clone: %v", err)
	}

	s, err := Status(cloneDir, testBranch)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}

	if !s.IsCloned {
		t.Error("expected IsCloned = true")
	}
}

func TestPullAlreadyUpToDate(t *testing.T) {
	_, workDir := createTestRepo(t)

	updated, err := Pull(workDir)
	if err != nil {
		t.Fatalf("Pull: %v", err)
	}
	if updated {
		t.Error("expected no updates")
	}
}
