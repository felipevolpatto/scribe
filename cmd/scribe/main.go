package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	cfg "github.com/felipevolpatto/scribe/internal/config"
	gitpkg "github.com/felipevolpatto/scribe/internal/git"
	md "github.com/felipevolpatto/scribe/internal/markdown"
	"github.com/felipevolpatto/scribe/internal/parser"
	"github.com/felipevolpatto/scribe/internal/tui"
	wf "github.com/felipevolpatto/scribe/internal/workflow"
)

func main() {
    root := &cobra.Command{Use: "scribe", Short: "Automate Conventional Commits changelogs"}

    var repoPath string
    var fromRef string

    newCmd := &cobra.Command{
        Use:   "new",
        Short: "Generate changelog for unreleased changes and print to stdout",
        RunE: func(cmd *cobra.Command, args []string) error {
            configuration, err := cfg.Load(repoPath)
            if err != nil {
                return err
            }

            ref := fromRef
            if ref == "" {
                tag, err := gitpkg.GetLatestTag(repoPath)
                if err == nil {
                    ref = tag
                }
            }

            commits, err := gitpkg.GetCommitsSince(repoPath, ref)
            if err != nil {
                return err
            }

            var parsedCommits []*parser.ParsedCommit
            for i := range commits {
                msg := strings.Split(commits[i].Message, "\n")[0]
                pc, err := parser.Parse(msg)
                if err != nil {
                    continue
                }
                if shouldIgnoreByScope(pc.Scope, configuration) {
                    continue
                }
                pc.Raw = &commits[i]
                parsedCommits = append(parsedCommits, pc)
            }

            curated, err := tui.Run(parsedCommits, configuration)
            if err != nil {
                return err
            }

            out, err := md.Render("Unreleased", curated, configuration)
            if err != nil {
                return err
            }
            fmt.Fprintln(os.Stdout, out)
            return nil
        },
    }
    newCmd.Flags().StringVar(&repoPath, "path", ".", "Path to the git repository")
    newCmd.Flags().StringVar(&fromRef, "from-ref", "", "Git ref to start from instead of the latest tag")

    var releaseRepoPath string
    var noInteractive bool

    releaseCmd := &cobra.Command{
        Use:   "release <version>",
        Short: "Generate changelog, prepend to CHANGELOG.md, commit and tag",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            version := args[0]
            configuration, err := cfg.Load(releaseRepoPath)
            if err != nil {
                return err
            }

            tag, _ := gitpkg.GetLatestTag(releaseRepoPath)
            commits, err := gitpkg.GetCommitsSince(releaseRepoPath, tag)
            if err != nil {
                return err
            }
            var parsedCommits []*parser.ParsedCommit
            for i := range commits {
                msg := strings.Split(commits[i].Message, "\n")[0]
                pc, err := parser.Parse(msg)
                if err != nil {
                    continue
                }
                if shouldIgnoreByScope(pc.Scope, configuration) {
                    continue
                }
                pc.Raw = &commits[i]
                parsedCommits = append(parsedCommits, pc)
            }

            curated := parsedCommits
            if !noInteractive {
                curated, err = tui.Run(parsedCommits, configuration)
                if err != nil {
                    return err
                }
            }

            content, err := md.Render(versionWithV(version), curated, configuration)
            if err != nil {
                return err
            }
            today := time.Now().Format("2006-01-02")
            header := fmt.Sprintf("## %s - %s\n\n", versionWithV(version), today)
            final := header + content + "\n"

            if err := wf.PrependToFile(filepathJoin(releaseRepoPath, "CHANGELOG.md"), final); err != nil {
                return err
            }
            if err := wf.CommitAndTag(releaseRepoPath, versionWithV(version)); err != nil {
                return err
            }
            return nil
        },
    }
    releaseCmd.Flags().StringVar(&releaseRepoPath, "path", ".", "Path to the git repository")
    releaseCmd.Flags().BoolVar(&noInteractive, "no-interactive", false, "Disable interactive TUI")

    root.AddCommand(newCmd, releaseCmd)

    if err := root.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func versionWithV(v string) string {
    if strings.HasPrefix(v, "v") {
        return v
    }
    return "v" + v
}

func filepathJoin(elem ...string) string {
    // Avoid importing path/filepath to keep main minimal and OS-agnostic
    return strings.Join(elem, string(os.PathSeparator))
}

func shouldIgnoreByScope(scope string, c *cfg.Config) bool {
    if scope == "" {
        return false
    }
    for _, s := range c.IgnoreScopes {
        if s == scope {
            return true
        }
    }
    return false
}


