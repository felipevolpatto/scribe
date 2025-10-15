package markdown

import (
	"bytes"
	"fmt"
	"text/template"

	cfg "github.com/felipevolpatto/scribe/internal/config"
	"github.com/felipevolpatto/scribe/internal/parser"
)

// Render takes the final curated commits, the config, and the new version string
// and returns the formatted changelog content as a string.
func Render(version string, commits []*parser.ParsedCommit, config *cfg.Config) (string, error) {
    type item struct {
        Title string
        Lines []string
    }
    byTitle := make([]item, 0, len(config.Sections))
    for _, section := range config.Sections {
        var lines []string
        for _, pc := range commits {
            // Route breaking changes to any section with empty types
            if pc.IsBreaking && len(section.Types) == 0 {
                hash := ""
                if pc.Raw != nil {
                    if len(pc.Raw.Hash) >= 7 {
                        hash = pc.Raw.Hash[:7]
                    } else {
                        hash = pc.Raw.Hash
                    }
                }
                suffix := ""
                if hash != "" {
                    suffix = fmt.Sprintf(" (%s)", hash)
                }
                lines = append(lines, "* "+pc.Description+suffix)
                continue
            }
            for _, t := range section.Types {
                if pc.Type == t {
                    hash := ""
                    if pc.Raw != nil {
                        if len(pc.Raw.Hash) >= 7 {
                            hash = pc.Raw.Hash[:7]
                        } else {
                            hash = pc.Raw.Hash
                        }
                    }
                    suffix := ""
                    if hash != "" {
                        suffix = fmt.Sprintf(" (%s)", hash)
                    }
                    lines = append(lines, "* "+pc.Description+suffix)
                    break
                }
            }
        }
        if len(lines) > 0 {
            byTitle = append(byTitle, item{Title: section.Title, Lines: lines})
        }
    }

    const tmpl = `{{- range . }}### {{ .Title }}
{{- range .Lines }}
{{ . }}
{{- end }}

{{ end -}}`

    t := template.Must(template.New("changelog").Parse(tmpl))
    var buf bytes.Buffer
    if err := t.Execute(&buf, byTitle); err != nil {
        return "", err
    }
    return buf.String(), nil
}


