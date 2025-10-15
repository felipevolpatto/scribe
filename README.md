# Scribe

Scribe is a CLI that automates `CHANGELOG.md` generation using Conventional Commits.

## Features

- Interactive TUI to curate commit messages (include, edit, re-categorize)
- Follows Conventional Commits; non-conforming commits are skipped
- Generates Markdown sections by type (feat, fix, etc.)
- Release workflow: prepend to `CHANGELOG.md`, commit, and create a git tag

## Install

```bash
go install github.com/felipevolpatto/scribe/cmd/scribe@latest
```

Or download prebuilt binaries from GitHub Releases (created on tag pushes).

### Requirements

- Go 1.22+

## Usage

- Generate changelog preview (does not modify files):
```bash
scribe new --path . [--from-ref <git_ref>]
```

- Create a release: prepend to `CHANGELOG.md`, commit, and tag:
```bash
scribe release 1.2.0 --path . [--no-interactive]
```
Notes:
- The git tag will be created as `v1.2.0` (Scribe prefixes with `v`).
- Use `--no-interactive` for CI or fully automated runs.

## Config (.scribe.yml)

```yaml
sections:
  - { title: "Breaking Changes", types: [] }
  - { title: "New Features", types: ["feat"] }
  - { title: "Bug Fixes", types: ["fix"] }
ignore_scopes: []
```

Behavior:
- Any section with empty `types` is treated as the Breaking Changes bucket; commits marked with `!` are routed there.
- `ignore_scopes` filters out commits whose scope matches any entry.

## TUI Keybindings

- Space: toggle include/exclude
- e: edit selected commit description
- c: change commit type (cycles through known types)
- Enter: confirm selection
- q / Esc: abort

## Development

```bash
go test ./...
```

## CI

- Tests run on every push via GitHub Actions (`.github/workflows/test.yml`).
- Creating a tag like `v1.2.0` triggers a Release build that uploads multi-arch binaries (`.github/workflows/release.yml`).

## License

MIT © Felipe Volpatto
