package tui

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	cfg "github.com/felipevolpatto/scribe/internal/config"
	"github.com/felipevolpatto/scribe/internal/parser"
)

var ErrAborted = errors.New("user aborted")

type mode int

const (
    modeNormal mode = iota
    modeEdit
)

type model struct {
    commits      []*parser.ParsedCommit
    include      []bool
    idx          int
    mode         mode
    editBuffer   string
    allowedTypes []string
    aborted      bool
}

func initialModel(commits []*parser.ParsedCommit, configuration *cfg.Config) model {
    include := make([]bool, len(commits))
    for i := range include {
        include[i] = true
    }
    // derive allowed types from config sections (unique, preserve order)
    seen := map[string]bool{}
    var types []string
    for _, s := range configuration.Sections {
        for _, t := range s.Types {
            if !seen[t] {
                seen[t] = true
                types = append(types, t)
            }
        }
    }
    // ensure some common types are present
    for _, t := range []string{"feat", "fix", "chore", "refactor", "docs", "test", "perf"} {
        if !seen[t] {
            seen[t] = true
            types = append(types, t)
        }
    }
    return model{commits: commits, include: include, allowedTypes: types}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        key := msg.String()
        if m.mode == modeEdit {
            switch msg.Type {
            case tea.KeyEnter:
                if m.idx >= 0 && m.idx < len(m.commits) {
                    m.commits[m.idx].Description = strings.TrimSpace(m.editBuffer)
                }
                m.mode = modeNormal
                m.editBuffer = ""
                return m, nil
            case tea.KeyEsc:
                m.mode = modeNormal
                m.editBuffer = ""
                return m, nil
            case tea.KeyBackspace:
                if len(m.editBuffer) > 0 {
                    m.editBuffer = m.editBuffer[:len(m.editBuffer)-1]
                }
                return m, nil
            default:
                if msg.Runes != nil {
                    m.editBuffer += string(msg.Runes)
                }
                return m, nil
            }
        }

        switch key {
        case "q", "esc":
            m.aborted = true
            return m, tea.Quit
        case "up", "k":
            if m.idx > 0 {
                m.idx--
            }
        case "down", "j":
            if m.idx < len(m.commits)-1 {
                m.idx++
            }
        case " ":
            if m.idx >= 0 && m.idx < len(m.include) {
                m.include[m.idx] = !m.include[m.idx]
            }
        case "e":
            if m.idx >= 0 && m.idx < len(m.commits) {
                m.mode = modeEdit
                m.editBuffer = m.commits[m.idx].Description
            }
        case "c":
            if m.idx >= 0 && m.idx < len(m.commits) {
                cur := m.commits[m.idx].Type
                next := nextType(cur, m.allowedTypes)
                m.commits[m.idx].Type = next
            }
        case "enter":
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m model) View() string {
    var b strings.Builder
    b.WriteString("Scribe - Select commits to include. Space: toggle, e: edit, c: retype, Enter: confirm, q: quit\n\n")
    for i, c := range m.commits {
        cursor := " "
        if i == m.idx {
            cursor = ">"
        }
        mark := "[x]"
        if !m.include[i] {
            mark = "[ ]"
        }
        scope := c.Scope
        if scope != "" {
            scope = fmt.Sprintf("(%s)", scope)
        }
        bang := ""
        if c.IsBreaking {
            bang = "!"
        }
        line := fmt.Sprintf("%s %s %s%s%s: %s\n", cursor, mark, c.Type, scope, bang, c.Description)
        b.WriteString(line)
    }
    if m.mode == modeEdit {
        b.WriteString("\nEditing description: " + m.editBuffer)
    }
    return b.String()
}

func nextType(cur string, types []string) string {
    if len(types) == 0 {
        return cur
    }
    for i, t := range types {
        if t == cur {
            return types[(i+1)%len(types)]
        }
    }
    return types[0]
}

// Run launches the interactive terminal UI.
// It takes the commits found by the parser and the loaded config.
// It returns the final, curated list of commits that the user has approved.
// Returns an error if the user aborts the session.
func Run(commits []*parser.ParsedCommit, configuration *cfg.Config) ([]*parser.ParsedCommit, error) {
    m := initialModel(commits, configuration)
    prog := tea.NewProgram(m)
    res, err := prog.Run()
    if err != nil {
        return nil, err
    }
    fm := res.(model)
    if fm.aborted {
        return nil, ErrAborted
    }
    var out []*parser.ParsedCommit
    for i, c := range fm.commits {
        if fm.include[i] {
            out = append(out, c)
        }
    }
    return out, nil
}



