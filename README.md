# MarkWatch - Markdown Preview Tool

This Go program allows you to preview Markdown files in a browser. It converts Markdown to HTML, sanitizes it for safe display, and serves it locally. The server auto-refreshes whenever the file changes, allowing you to see live updates.

## Features

- Converts Markdown to HTML with a user-defined template
- Real-time auto-refresh in the browser when file changes
- Safe HTML output via sanitization
- CLI flags for customization

## Requirements

- Go 1.16 or later
- `fsnotify` for file-watching

## Installation

1. Clone this repository:

   ```
   git clone <repository-url>
   cd <repository-folder>
   ```

2. Install dependencies (if necessary):

   ```
   go mod tidy
   ```

3. Build the application:

   ```
   go build -o markwatch
   ```

## Usage

```bash
./md-preview -file <path/to/markdown-file.md> [-t <path/to/template.html>] [-s]
```

### Command-line Flags

- `-file <path>`: Specify the Markdown file to preview. **(Required)**
- `-t <path>`: Optional custom HTML template for rendering.
- `-s`: Skip auto-preview in the browser on startup.

## Example

```bash
./md-preview -file README.md -t templates/custom.html
```

## How It Works

1. **Markdown Parsing**: Uses the `blackfriday` package to parse Markdown into HTML.
2. **HTML Sanitization**: `bluemonday` sanitizes HTML to prevent potential XSS attacks.
3. **Auto-reload**: Monitors the file for changes, automatically refreshing the preview.
4. **Server Shutdown**: Gracefully stops the server on receiving system interrupts.

## File Structure

- `templates/default.html`: Default HTML template
- `main.go`: Core logic of the application

## Troubleshooting

- **Port Conflict**: The server runs on port `9090` by default. If this port is in use, ensure no other services are running on it.
- **Template Errors**: Custom templates must define `{{ .Title }}` and `{{ .Body }}` placeholders.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.

Enjoy live Markdown previewing!
