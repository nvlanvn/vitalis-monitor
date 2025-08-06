# Vitalis Monitor

A modern terminal user interface for the `netstat` command, built with Go and Bubble Tea.

## Features

- 🎨 Beautiful tabbed interface for different netstat outputs
- ⌨️ Keyboard navigation (arrow keys or vim-style h/l)
- 📊 Interactive tables with scrolling support
- 🎯 Smart column width calculation
- 🔍 Automatic parsing of various netstat formats

## Installation

```bash
go install github.com/nvlanvn/vitalis-monitor@latest
```

Or build from source:

```bash
git clone https://github.com/nvlanvn/vitalis-monitor.git
cd vitalis-monitor
go build -o vitalis-monitor
```

## Usage

Run with the same arguments you would pass to netstat:

```bash
# Show all connections
vitalis-monitor -a

# Show listening ports
vitalis-monitor -l

# Show network statistics
vitalis-monitor -s

# Show routing table
vitalis-monitor -r
```

### Keyboard Shortcuts

- `←` / `h` - Navigate to previous tab
- `→` / `l` - Navigate to next tab
- `↑` / `↓` - Navigate table rows
- `q` / `Ctrl+C` / `ESC` - Quit

## Project Structure

```
vitalis-monitor/
├── main.go                    # Application entry point
├── go.mod                     # Go module definition
├── internal/
│   ├── netstat/              # Netstat command execution
│   │   └── runner.go
│   ├── parser/               # Output parsing logic
│   │   └── parser.go
│   ├── styles/               # UI styling definitions
│   │   └── styles.go
│   └── ui/                   # Terminal UI components
│       ├── model.go          # Main UI model
│       └── table.go          # Table building logic
└── README.md
```

## Architecture

The application follows a modular architecture with clear separation of concerns:

1. **Netstat Runner**: Executes the netstat command and provides output stream
2. **Parser**: Intelligently parses netstat output into structured sections
3. **UI Model**: Manages application state and user interactions
4. **Table Builder**: Converts parsed data into interactive tables
5. **Styles**: Centralized styling configuration

## Development

### Prerequisites

- Go 1.21 or higher
- netstat command available in PATH

### Building

```bash
go build -o vitalis-monitor
```

### Running Tests

```bash
go test ./...
```

### Adding New Features

1. **New Netstat Formats**: Add patterns to `parser.knownSections` and header mappings to `parser.headerMappings`
2. **New Styles**: Modify the `styles` package
3. **New UI Features**: Extend the `ui.Model` with new functionality

## Contributing

Pull requests are welcome! Please ensure:

1. Code follows Go conventions
2. All tests pass
3. New features include tests
4. Documentation is updated

## License

MIT License - see LICENSE file for details
