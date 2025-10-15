package parser

import "testing"

func TestParse_ValidMessages(t *testing.T) {
    tests := []struct {
        msg        string
        typ        string
        scope      string
        desc       string
        breaking   bool
    }{
        {"feat: add login", "feat", "", "add login", false},
        {"fix(api): correct bug", "fix", "api", "correct bug", false},
        {"feat(ui)!: make button primary", "feat", "ui", "make button primary", true},
        {"chore: bump deps", "chore", "", "bump deps", false},
    }

    for _, tt := range tests {
        parsed, err := Parse(tt.msg)
        if err != nil {
            t.Fatalf("unexpected error for %q: %v", tt.msg, err)
        }
        if parsed.Type != tt.typ || parsed.Scope != tt.scope || parsed.Description != tt.desc || parsed.IsBreaking != tt.breaking {
            t.Fatalf("parsed mismatch for %q: got %+v", tt.msg, parsed)
        }
    }
}

func TestParse_InvalidMessage(t *testing.T) {
    if _, err := Parse("this is not conventional"); err == nil {
        t.Fatal("expected error for invalid message, got nil")
    }
}


