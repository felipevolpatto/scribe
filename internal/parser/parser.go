package parser

import (
	"errors"
	"regexp"

	gitpkg "github.com/felipevolpatto/scribe/internal/git"
)

// ParsedCommit is the structured representation of a Conventional Commit.
type ParsedCommit struct {
    Type        string
    Scope       string
    Description string
    IsBreaking  bool
    Raw         *gitpkg.RawCommit
}

var (
    // Matches: type(scope)!: description
    // Groups: 1=type 2=scope (optional) 3=! (optional) 4=description
    conventionalRe = regexp.MustCompile(`^(\w+)(?:\(([\w\/-]+)\))?(!)?:\s+(.+)$`)
)

// Parse takes a raw commit message and returns a structured ParsedCommit.
// Returns an error if the message does not conform to the spec.
func Parse(message string) (*ParsedCommit, error) {
    if message == "" {
        return nil, errors.New("empty commit message")
    }

    matches := conventionalRe.FindStringSubmatch(message)
    if matches == nil {
        return nil, errors.New("commit message does not follow Conventional Commits")
    }

    parsed := &ParsedCommit{
        Type:        matches[1],
        Scope:       matches[2],
        Description: matches[4],
        IsBreaking:  matches[3] == "!",
    }
    return parsed, nil
}


