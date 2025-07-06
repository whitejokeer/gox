# GOX CLI Commands

The GOX CLI provides a complete set of commands for developing with the GOX framework.

## Commands

### `gox init [name]`
Create a new GOX project with the basic structure and configuration.

**Examples:**
```bash
gox init my-app          # Create project in ./my-app
gox init .               # Initialize in current directory
```

**Creates:**
- Project directory structure (`src/`, `static/`, `dist/`)
- Configuration file (`gox.toml`)
- Example component (`src/components/welcome.gox`)
- Basic Go main file
- `.gitignore` file

### `gox create [type] [name]`
Generate new components, pages, or services with boilerplate code.

**Available types:**
- `component` - Create a new component
- `page` - Create a new page component
- `service` - Create a new service

**Examples:**
```bash
gox create component button       # Creates src/components/button.gox
gox create page about            # Creates src/pages/about.gox
gox create service auth          # Creates src/services/auth.go
```

### `gox dev`
Start development server with hot reload functionality.

**Flags:**
- `--port, -p` - Port to run the server on (default: 3000)
- `--host, -H` - Host to bind the server to (default: localhost)

**Examples:**
```bash
gox dev                          # Start on default port (3000)
gox dev --port 8080              # Start on custom port
gox dev --host 0.0.0.0           # Bind to all interfaces
```

### `gox build`
Build .gox files into Go components for production.

**Flags:**
- `--output, -o` - Output directory for built files (default: ./dist)

**Examples:**
```bash
gox build                        # Build entire project
gox build --output ./custom-dist # Custom output directory
```

### `gox watch`
Watch .gox files for changes and rebuild automatically.

**Examples:**
```bash
gox watch                        # Watch all .gox files
```

## Configuration

GOX uses a `gox.toml` configuration file:

```toml
[project]
name = "my-app"
version = "0.1.0"

[server]
port = 3000
host = "localhost"

[build]
output = "./dist"
minify = true

[dev]
hot_reload = true
watch_paths = ["src/", "static/"]
```

## Global Flags

- `--config` - Config file path (default: ./gox.toml)
- `--no-color` - Disable colored output
- `--help, -h` - Show help
- `--version, -v` - Show version

## Environment Variables

Configuration can be overridden with environment variables using the `GOX_` prefix:

```bash
export GOX_SERVER_PORT=8080
export GOX_BUILD_OUTPUT=./custom-dist
```

## Getting Help

Use `gox [command] --help` for detailed help on any command.