# Skill Router

A local web application to manage Claude Code skills and plugins.

![Main View](docs/images/main-view.png)

## Features

- **View all skills** from `~/.claude/commands/` and installed plugins
- **Enable/disable skills** individually or by plugin group
- **Delete skills** (user skills only)
- **Upload .md skill files** via drag-and-drop or file picker
- **Install skills from GitHub** repositories
- **Multi-language support** - English and Chinese with auto-detection

## Installation

Download the binary for your platform from [Releases](https://github.com/anthropics/skill-router/releases).

Or build from source:

```bash
make build
```

## Usage

```bash
./skill-router
```

This starts the server and opens your browser to http://localhost:9527

### Adding Skills

Click the **+ Add** button to:

1. **Upload a file** - Drag and drop or select a `.md` skill file
2. **Install from GitHub** - Enter a repository URL to install all skills from `.claude/commands/`

![Add Skill Modal](docs/images/add-skill-modal.png)

### Language

The interface automatically detects your browser language. Click the language toggle (EN/中) in the header to switch manually.

## Development

```bash
# Terminal 1: Run Go backend
go run .

# Terminal 2: Run Vue dev server
cd web && npm run dev
```

Then open http://localhost:5173 for hot-reload development.

### Project Structure

```
.
├── main.go                 # Entry point, HTTP server
├── internal/
│   ├── handler/            # HTTP handlers
│   ├── service/            # Business logic
│   └── config/             # Configuration management
├── web/                    # Vue 3 frontend
│   ├── src/
│   │   ├── components/     # Vue components
│   │   ├── i18n/           # Internationalization
│   │   └── api/            # API client
│   └── dist/               # Built frontend (embedded)
└── docs/
    └── images/             # Screenshots
```

## License

MIT
