package workflow

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// PrependToFile adds the new changelog content to the top of CHANGELOG.md.
// It should handle creating the file if it doesn't exist.
func PrependToFile(filePath string, content string) error {
    var existing []byte
    if b, err := os.ReadFile(filePath); err == nil {
        existing = b
    }
    var buf bytes.Buffer
    buf.WriteString(content)
    buf.Write(existing)
    return os.WriteFile(filePath, buf.Bytes(), 0o644)
}

// CommitAndTag executes 'git add CHANGELOG.md', 'git commit', and 'git tag'.
func CommitAndTag(repoPath, version string) error {
    cmds := [][]string{
        {"git", "add", "CHANGELOG.md"},
        {"git", "commit", "-m", fmt.Sprintf("chore(release): %s", version)},
        {"git", "tag", version},
    }
    for _, c := range cmds {
        cmd := exec.Command(c[0], c[1:]...)
        cmd.Dir = repoPath
        if out, err := cmd.CombinedOutput(); err != nil {
            return fmt.Errorf("%s failed: %v: %s", c, err, string(out))
        }
    }
    return nil
}


