# Skill Router

A local web application to manage Claude Code skills.

## Features

- View all skills from `~/.claude/commands/`
- Enable/disable skills (moves to `~/.claude/skills-disabled/`)
- Delete skills
- Upload .md skill files
- Install skills from GitHub repositories

## Installation

Download the binary for your platform from [Releases](https://github.com/wind/skill-router/releases).

Or build from source:

```bash
make build
```

## Usage

```bash
./skill-router
```

This starts the server and opens your browser to http://localhost:9527

## Development

```bash
# Terminal 1: Run Go backend
go run .

# Terminal 2: Run Vue dev server
cd web && npm run dev
```

## License

MIT
