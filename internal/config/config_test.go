package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_DefaultWhenMissing(t *testing.T) {
    dir := t.TempDir()
    c, err := Load(dir)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(c.Sections) == 0 {
        t.Fatalf("expected default sections")
    }
}

func TestLoad_FromFile(t *testing.T) {
    dir := t.TempDir()
    content := []byte("sections:\n  - { title: 'Custom', types: ['feat','fix'] }\nignore_scopes: ['docs']\n")
    if err := os.WriteFile(filepath.Join(dir, ".scribe.yml"), content, 0o644); err != nil {
        t.Fatalf("write config: %v", err)
    }
    c, err := Load(dir)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(c.Sections) != 1 || c.Sections[0].Title != "Custom" {
        t.Fatalf("expected custom section, got %+v", c.Sections)
    }
    if len(c.IgnoreScopes) != 1 || c.IgnoreScopes[0] != "docs" {
        t.Fatalf("expected ignore_scopes ['docs'], got %+v", c.IgnoreScopes)
    }
}


