# wsm — Workspace Manager

A single static binary that manages multi-repo workspaces. Define your
repositories and projects in `workspace.json`, then use `wsm` to clone,
monitor, and update them all.

## Install

```bash
# From source
go install github.com/thasso/wsm/cmd/wsm@latest

# Or build locally
make build      # produces bin/wsm
make install    # installs to $GOPATH/bin
```

## Quick Start

```bash
# Create a new workspace
wsm init --org myorg --jira-url https://jira.example.com

# Add repositories
wsm add-repo --name my-service --repo my-service --branch main --jira-key SVC

# Clone all repos and set up .gitignore
wsm setup

# Check status across all repos
wsm status

# Pull latest changes everywhere
wsm pull
```

## Commands

| Command | Description |
|---------|-------------|
| `wsm setup` | Clone missing repos, init submodules, update `.gitignore` |
| `wsm status` | Show branch, clean/dirty state, and ahead/behind for all repos |
| `wsm pull` | Pull latest changes in all repos on their current branch |
| `wsm init` | Create a new `workspace.json` manifest |
| `wsm add-repo` | Add a repository entry to `workspace.json` |

Run `wsm --help` or `wsm <command> --help` for full details on any command.

### Machine-readable output

```bash
wsm status --json    # JSON array of repo status objects
```

## workspace.json

The manifest has two sections:

- **`repos`** — Git repositories cloned into the workspace
- **`projects`** — Jira-only projects with no associated repository

### Top-level fields

| Field | Description |
|-------|-------------|
| `org` | GitHub organization (used for clone URLs) |
| `jira_url` | Base URL of the Jira instance |

### Repo entries (`repos[]`)

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | Local directory name |
| `repo` | yes | GitHub repository name under `org` |
| `branch` | yes | Default branch to clone |
| `jira_key` | no | Jira project key |
| `description` | no | Brief description |

### Project entries (`projects[]`)

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | Human-readable project name |
| `jira_key` | yes | Jira project key |
| `description` | no | Brief description |

### Example

```json
{
  "org": "thasso",
  "jira_url": "https://jira.example.com",
  "repos": [
    {
      "name": "backend-api",
      "repo": "backend-api",
      "branch": "main",
      "jira_key": "API",
      "description": "Core backend REST API"
    },
    {
      "name": "web-frontend",
      "repo": "web-frontend",
      "branch": "main",
      "jira_key": "WEB",
      "description": "React web frontend"
    }
  ],
  "projects": [
    {
      "name": "Infrastructure",
      "jira_key": "INFRA",
      "description": "Infra and DevOps tracking project"
    }
  ]
}
```

## Working with Sub-Repos

Each sub-repo is a normal git clone. Just `cd` into it and use git as usual:

```bash
cd backend-api
git checkout feature/my-branch
# ... make changes, commit, push ...
```

## Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-c, --config` | `workspace.json` | Path to workspace manifest |
| `-v, --version` | | Print version |

## For AI Agents

Run `wsm --help` to discover all workspace management commands. All commands
use POSIX-compliant flags and produce structured output (`--json` where
available). No interactive prompts — everything is driven by flags.
