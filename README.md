# ralph-cli

Autonomous PRD-driven agent loop. Reads stories from a `prd.json` file and executes them one at a time using AI coding agents (OpenCode, Claude Code, Codex, etc.). Files and git serve as memory for fresh iterations and persistent state.

## Installation

### Via supercli (recommended)

```bash
npm install -g superacli
sc plugins install ralph-cli
```

Then use through supercli:

```bash
sc ralph-cli init run "My Project" --prd ./prd.json
sc ralph-cli run run --prd ./prd.json
sc ralph-cli status run --prd ./prd.json
```

### From GitHub Releases

```bash
curl -LO https://github.com/javimosch/ralph-cli/releases/latest/download/ralph-linux-amd64
chmod +x ralph-linux-amd64
sudo mv ralph-linux-amd64 /usr/local/bin/ralph
```

### Build from source

```bash
git clone https://github.com/javimosch/ralph-cli.git
cd ralph-cli
go build -o ralph .
sudo mv ralph /usr/local/bin/
```

## Quick Start

```bash
# Initialize a PRD
ralph init "My Feature" --prd ./prd.json

# Check status
ralph status --prd ./prd.json

# Dry run (preview prompts without executing)
ralph run --prd ./prd.json --dry-run

# Run the agent loop
ralph run --prd ./prd.json
```

## PRD Format

```json
{
  "name": "My Feature",
  "branchName": "ralph/my-feature",
  "description": "What this PRD is about",
  "userStories": [
    {
      "id": "US-001",
      "title": "Add database schema",
      "description": "As a developer, I need...",
      "acceptanceCriteria": ["Column x exists", "Migration runs"],
      "priority": 1,
      "passes": false,
      "dependsOn": []
    }
  ]
}
```

## Commands

| Command | Description |
|---------|-------------|
| `ralph run --prd <file>` | Execute the agent loop |
| `ralph status --prd <file>` | Show progress |
| `ralph story next --prd <file>` | Show next available story |
| `ralph story prompt --prd <file>` | Print agent prompt for a story |
| `ralph init <name> --prd <file>` | Scaffold a new PRD |

## Flags

| Flag | Description | Default |
|------|-------------|--------|
| `--prd` | Path to PRD JSON file | required |
| `--agent` | Agent CLI command | auto-detect |
| `--dry-run` | Print prompts without executing | false |
| `--timeout` | Max seconds per agent execution | 300 |

## supercli Usage

Primary usage is through supercli (install via `npm install -g superacli`):

```bash
# Install the plugin
sc plugins install ralph-cli

# Init a PRD
sc ralph-cli init run "Feature name" --prd ./prd.json

# Run the loop
sc ralph-cli run run --prd ./prd.json

# Check status
sc ralph-cli status run --prd ./prd.json

# Next story
sc ralph-cli story next --prd ./prd.json

# Agent prompt
sc ralph-cli prompt run --prd ./prd.json
```

## License

MIT — Copyright (c) 2025 Javier Leandro Arancibia
