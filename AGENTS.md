# AGENTS.md

`wsm` is a Go CLI tool that manages multi-repo workspaces. It reads a
`workspace.json` manifest and can clone, monitor, and update all listed
repositories. Run `wsm --help` for the full command reference.

## Quality Gates

Before submitting any changes, run these in order. All must pass.

```bash
make fmt        # Format all Go source files
make lint       # Run golangci-lint (zero issues required)
make build      # Compile the binary to bin/wsm
make test       # Run all tests
```

If `make fmt` changes any files, stage them before committing. CI will
reject unformatted code, lint errors, build failures, or test failures.

## Releasing

Tag with `git tag vX.Y.Z` and push with `git push --tags`. The release
workflow cross-compiles for linux (amd64/arm64), macOS (amd64/arm64),
and windows (amd64), then creates a GitHub Release with all binaries
and checksums attached.
