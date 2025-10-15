package workflow

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestPrependToFile(t *testing.T) {
    dir := t.TempDir()
    file := filepath.Join(dir, "CHANGELOG.md")

    if err := PrependToFile(file, "A\n"); err != nil {
        t.Fatalf("first prepend: %v", err)
    }
    if err := PrependToFile(file, "B\n"); err != nil {
        t.Fatalf("second prepend: %v", err)
    }
    b, _ := os.ReadFile(file)
    if string(b) != "B\nA\n" {
        t.Fatalf("unexpected content: %q", string(b))
    }
}

func TestCommitAndTag_Smoke(t *testing.T) {
    dir := t.TempDir()
    // init empty git repo
    run := func(name string, args ...string) {
        cmd := exec.Command(name, args...)
        cmd.Dir = dir
        if out, err := cmd.CombinedOutput(); err != nil {
            t.Fatalf("%s %v failed: %v: %s", name, args, err, string(out))
        }
    }
    run("git", "init")
    run("git", "config", "user.email", "test@example.com")
    run("git", "config", "user.name", "Test User")

    // create a file to commit
    if err := os.WriteFile(filepath.Join(dir, "CHANGELOG.md"), []byte("init\n"), 0o644); err != nil {
        t.Fatal(err)
    }

    if err := CommitAndTag(dir, "v0.0.1"); err != nil {
        t.Fatalf("CommitAndTag: %v", err)
    }
}


