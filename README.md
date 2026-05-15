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
curl -LO https://github.com/javimosch/ralph-cli/releases/latest/download/ralph-cli-linux-amd64
chmod +x ralph-cli-linux-amd64
sudo mv ralph-cli-linux-amd64 /usr/local/bin/ralph-cli
```

### Build from source

```bash
git clone https://github.com/javimosch/ralph-cli.git
cd ralph-cli
go build -o ralph-cli .
sudo mv ralph-cli /usr/local/bin/ralph-cli
```

## Quick Start

```bash
# 1. Initialize a PRD
ralph-cli init "My Feature" --prd ./tasks/prd.json

# 2. Check what needs to be done
ralph-cli status --prd ./tasks/prd.json
ralph-cli story next --prd ./tasks/prd.json

# 3. Preview the agent prompt
ralph-cli story prompt --prd ./tasks/prd.json

# 4. Run the agent loop (dry-run first to preview)
ralph-cli run --prd ./tasks/prd.json --dry-run
ralph-cli run --prd ./tasks/prd.json
```

## Commands

| Command | Description |
|---------|-------------|
| `ralph-cli init <name> [description]` | Scaffold a new PRD JSON file |
| `ralph-cli run --prd <file>` | Execute the agent loop (story by story) |
| `ralph-cli status --prd <file>` | Show progress, next story, blocked stories |
| `ralph-cli story next --prd <file>` | Print the next available story |
| `ralph-cli story prompt --prd <file> [--story <id>]` | Generate agent prompt for a story |

## Options

| Flag | Description |
|------|-------------|
| `--prd <path>` | Path to PRD JSON file (required for most commands) |
| `--agent <cmd>` | Agent CLI to use (default: auto-detect opencode, claude, codex) |
| `--dry-run` | Print prompts without executing |
| `--timeout <sec>` | Max seconds per agent execution (default: no timeout) |

## PRD Format

```json
{
  "name": "My Feature",
  "branchName": "my-feature",
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

## How it Works

1. **Define stories** in a `prd.json` with priorities, dependencies, and acceptance criteria
2. **Run the loop** — `ralph-cli` selects the next unblocked story, builds a prompt, and spawns an agent
3. **Agent implements** — the agent works autonomously and commits changes
4. **Repeat** — `ralph-cli` marks the story as passing and moves to the next

## License

MIT
