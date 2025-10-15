package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGit_TagAndCommitsSince(t *testing.T) {
    dir := t.TempDir()
    run := func(args ...string) {
        cmd := exec.Command("git", args...)
        cmd.Dir = dir
        if out, err := cmd.CombinedOutput(); err != nil {
            t.Fatalf("git %v failed: %v: %s", args, err, string(out))
        }
    }
    run("init")
    run("config", "user.email", "test@example.com")
    run("config", "user.name", "Test User")

    // initial commit
    if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0o644); err != nil {
        t.Fatal(err)
    }
    run("add", "a.txt")
    run("commit", "-m", "chore: init")
    run("tag", "v0.1.0")

    // new commits
    if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0o644); err != nil {
        t.Fatal(err)
    }
    run("add", "b.txt")
    run("commit", "-m", "feat: add b")
    run("tag", "v0.1.1")
    if err := os.WriteFile(filepath.Join(dir, "c.txt"), []byte("c"), 0o644); err != nil {
        t.Fatal(err)
    }
    run("add", "c.txt")
    run("commit", "-m", "fix: add c")

    tag, err := GetLatestTag(dir)
    if err != nil {
        t.Fatalf("GetLatestTag: %v", err)
    }
    if tag != "v0.1.1" {
        t.Fatalf("expected latest tag v0.1.1, got %s", tag)
    }

    commits, err := GetCommitsSince(dir, "v0.1.0")
    if err != nil {
        t.Fatalf("GetCommitsSince: %v", err)
    }
    if len(commits) < 2 {
        t.Fatalf("expected >=2 commits since v0.1.0, got %d", len(commits))
    }
}


