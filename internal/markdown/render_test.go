package markdown

import (
	"strings"
	"testing"

	cfg "github.com/felipevolpatto/scribe/internal/config"
	gitpkg "github.com/felipevolpatto/scribe/internal/git"
	"github.com/felipevolpatto/scribe/internal/parser"
)

func TestRender_BasicSections(t *testing.T) {
    config := cfg.Default()
    commits := []*parser.ParsedCommit{
        {Type: "feat", Description: "add login", Raw: &gitpkg.RawCommit{Hash: "abcdef1"}},
        {Type: "fix", Description: "correct bug", Raw: &gitpkg.RawCommit{Hash: "1234567"}},
        {Type: "feat", Description: "breaking api", IsBreaking: true, Raw: &gitpkg.RawCommit{Hash: "7654321"}},
    }
    out, err := Render("v1.0.0", commits, config)
    if err != nil {
        t.Fatalf("render error: %v", err)
    }
    if !strings.Contains(out, "### Breaking Changes") {
        t.Fatalf("missing breaking changes section: %q", out)
    }
    if !strings.Contains(out, "### New Features") || !strings.Contains(out, "add login (abcdef1)") {
        t.Fatalf("features not rendered correctly: %q", out)
    }
    if !strings.Contains(out, "### Bug Fixes") || !strings.Contains(out, "correct bug (1234567)") {
        t.Fatalf("bug fixes not rendered correctly: %q", out)
    }
}

func TestRender_OmitEmptyHashParentheses(t *testing.T) {
    config := cfg.Default()
    commits := []*parser.ParsedCommit{
        {Type: "feat", Description: "no hash"},
    }
    out, err := Render("v1.0.0", commits, config)
    if err != nil {
        t.Fatalf("render error: %v", err)
    }
    if strings.Contains(out, "()") {
        t.Fatalf("should not render empty parentheses: %q", out)
    }
}

func TestRender_BreakingSectionEmptyTypes(t *testing.T) {
    // Custom config: first section has empty types and should catch breaking changes
    config := &cfg.Config{Sections: []cfg.Section{
        {Title: "Custom Breaking", Types: []string{}},
        {Title: "Features", Types: []string{"feat"}},
    }}
    commits := []*parser.ParsedCommit{
        {Type: "feat", Description: "add a", Raw: &gitpkg.RawCommit{Hash: "abcdef1"}},
        {Type: "feat", Description: "breaking!", IsBreaking: true, Raw: &gitpkg.RawCommit{Hash: "1234567"}},
    }
    out, err := Render("v1.0.0", commits, config)
    if err != nil {
        t.Fatalf("render error: %v", err)
    }
    if !strings.Contains(out, "### Custom Breaking") {
        t.Fatalf("missing custom breaking section: %q", out)
    }
    if !strings.Contains(out, "breaking! (1234567)") {
        t.Fatalf("missing breaking commit under custom section: %q", out)
    }
}


