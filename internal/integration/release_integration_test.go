package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestReleaseCommand_GeneratesChangelogAndTag(t *testing.T) {
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

    if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0o644); err != nil {
        t.Fatal(err)
    }
    run("add", "a.txt")
    run("commit", "-m", "chore: init")
    run("tag", "v0.1.0")

    if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0o644); err != nil {
        t.Fatal(err)
    }
    run("add", "b.txt")
    run("commit", "-m", "feat: add feature")
    if err := os.WriteFile(filepath.Join(dir, "c.txt"), []byte("c"), 0o644); err != nil {
        t.Fatal(err)
    }
    run("add", "c.txt")
    run("commit", "-m", "fix: bug fix")

    // module root is the project root which contains go.mod; this test is in internal/integration
    moduleRoot := filepath.Join("..", "..")
    if runtime.GOOS == "windows" {
        // no-op; relative works on windows too, kept for clarity
    }
    cmd := exec.Command("go", "run", "./cmd/scribe", "release", "0.2.0", "--no-interactive", "--path", dir)
    cmd.Dir = moduleRoot
    if out, err := cmd.CombinedOutput(); err != nil {
        t.Fatalf("release run failed: %v: %s", err, string(out))
    }

    if _, err := os.Stat(filepath.Join(dir, "CHANGELOG.md")); err != nil {
        t.Fatalf("CHANGELOG.md missing: %v", err)
    }
    // verify tag exists
    tagList := exec.Command("git", "tag", "--list", "v0.2.0")
    tagList.Dir = dir
    if out, err := tagList.CombinedOutput(); err != nil || string(out) == "" {
        t.Fatalf("expected tag v0.2.0, got err=%v out=%q", err, string(out))
    }
}


