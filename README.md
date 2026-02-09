# tviewyaml

[![CI](https://github.com/cassdeckard/tviewyaml/actions/workflows/ci.yml/badge.svg)](https://github.com/cassdeckard/tviewyaml/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/cassdeckard/tviewyaml/main/coverage.json)](https://github.com/cassdeckard/tviewyaml/actions/workflows/ci.yml)

A Go package for building [tview](https://github.com/rivo/tview) terminal UI applications from YAML configuration files.

## Overview

`tviewyaml` allows you to define your terminal UI layouts, widgets, and interactions in YAML files instead of writing Go code. This makes it easier to:

- Design and prototype terminal UIs quickly
- Separate UI structure from application logic
- Modify layouts without recompiling
- Create reusable UI components

## Installation

```bash
go get github.com/cassdeckard/tviewyaml
```

## Features

- **Declarative UI**: Define your entire UI structure in YAML
- **Rich Widget Support**: Lists, Forms, Tables, Trees, Grids, Flex layouts, and more
- **Template Functions**: Built-in functions for page navigation, modals, and custom callbacks
- **Color Support**: Easy color configuration with named colors
- **Validation**: Built-in configuration validation

## Supported Widgets

- Box
- TextView
- Button
- List
- Flex (horizontal/vertical layouts)
- Form (with InputField, Checkbox, Dropdown, Button)
- Table
- TreeView
- Grid
- Modal
- Pages

## Quick Start

### 1. Create an Application Configuration (`app.yaml`)

```yaml
application:
  name: "My TView App"
  enableMouse: true
  globalKeyBindings:
    - key: "Escape"
      action: '{{ switchToPage "main" }}'
    - key: "Ctrl+Q"
      action: '{{ stopApp }}'
  root:
    type: pages
    pages:
      - name: main
        ref: main.yaml
      - name: settings
        ref: settings.yaml
```

### 2. Create a Page Configuration (`main.yaml`)

```yaml
type: list
title: Main Menu
border: true
titleAlign: center
listItems:
  - mainText: Start Application
    secondaryText: Launch the main app
    shortcut: "s"
    onSelected: "{{ showSimpleModal \"Application Started!\" \"OK\" }}"
  
  - mainText: Settings
    secondaryText: Configure application
    shortcut: "c"
    onSelected: "{{ switchToPage \"settings\" }}"
  
  - mainText: Exit
    secondaryText: Quit the application
    shortcut: "q"
    onSelected: "{{ stopApp }}"
```

### 3. Use in Your Go Application

```go
package main

import (
    "log"
    "github.com/cassdeckard/tviewyaml"
)

func main() {
    // Create app from YAML config directory using the Builder pattern
    app, err := tviewyaml.NewAppBuilder("./config").
        Build()
    if err != nil {
        log.Fatal(err)
    }

    // Run the application
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## Configuration Structure

The root configuration file defines application-level settings and the root view:

- **`application`**: Top-level application configuration
  - **`name`**: Application name (optional)
  - **`enableMouse`**: Enable mouse support (optional, defaults to true)
  - **`globalKeyBindings`**: Array of global keyboard shortcuts
    - **`key`**: Key string (e.g., "Escape", "Ctrl+Q", "F1")
    - **`action`**: Template expression to execute
  - **`root`**: The root view definition (currently must be type "pages")
    - **`type`**: View type (currently only "pages" supported)
    - **`pages`**: Array of page references

## Examples

The [`example/`](example/) directory contains a comprehensive demonstration application showcasing all widget types and features. This is the best way to learn how to use tviewyaml.

### Running the Examples

```bash
cd example
go run main.go
```

### What's Included

The example application demonstrates:

- **Layouts**: Flex (horizontal/vertical), Grid (responsive)
- **Widgets**: Box, Button, Checkbox, Dropdown, Form, InputField, List, Modal, Table, TextView, TreeView
- **Features**: Borders, titles, colors, text alignment, dynamic colors, regions
- **Callbacks**: Page navigation, modals, form submissions, app control
- **Patterns**: Nested layouts, complex forms, data tables, hierarchical trees

### Navigation

- Use arrow keys and Tab to navigate between items
- Press shortcut keys (shown in the main menu) to jump to specific demos
- Press ESC to return to the main menu from any page

See the [example README](example/README.md) for detailed information about each demonstration.

## Template Functions

Built-in template functions for callbacks:

- `switchToPage "pageName"` - Navigate to a different page
- `removePage "pageName"` - Remove a page from the stack
- `stopApp` - Exit the application
- `showSimpleModal "text" "button1" "button2"` - Show a modal dialog
- `noop` - No operation (placeholder callback)

### Custom Template Functions

You can register custom template functions using the Builder API. Each function is defined by:

- **Name**: String identifier used in templates (e.g., `"myCustomFunc"`)
- **MinArgs**: Minimum number of arguments (non-negative integer)
- **MaxArgs**: Maximum number of arguments (`nil` for unlimited/variadic)
- **Validator**: Optional validation function (called after argument count is validated)
- **Handler**: Function that executes the template logic

#### Example: Fixed Arguments

```go
package main

import (
    "log"
    "github.com/cassdeckard/tviewyaml"
    "github.com/cassdeckard/tviewyaml/template"
)

func main() {
    // Helper to create *int for maxArgs
    intPtr := func(i int) *int { return &i }
    
    app, err := tviewyaml.NewAppBuilder("./config").
        WithTemplateFunction("logMessage", 1, intPtr(1), nil,
            func(ctx *template.Context, message string) {
                log.Printf("Custom log: %s", message)
            },
        ).
        WithTemplateFunction("switchAndLog", 2, intPtr(2), nil,
            func(ctx *template.Context, pageName, logMsg string) {
                log.Printf("Switching to %s: %s", pageName, logMsg)
                ctx.Pages.SwitchToPage(pageName)
            },
        ).
        Build()
    if err != nil {
        log.Fatal(err)
    }

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

Then use in your YAML:

```yaml
onSelected: "{{ logMessage \"Button clicked!\" }}"
onSelected: "{{ switchAndLog \"settings\" \"User opened settings\" }}"
```

#### Example: Variadic Arguments

```go
app, err := tviewyaml.NewAppBuilder("./config").
    WithTemplateFunction("multiLog", 1, nil, nil,
        func(ctx *template.Context, args []string) {
            // Variadic - accepts 1 or more arguments
            for i, arg := range args {
                log.Printf("Arg %d: %s", i, arg)
            }
        },
    ).
    Build()
```

Use in YAML:

```yaml
onSelected: "{{ multiLog \"first\" \"second\" \"third\" }}"
```

#### Example: With Validator

```go
app, err := tviewyaml.NewAppBuilder("./config").
    WithTemplateFunction("switchToExistingPage", 1, intPtr(1),
        func(ctx *template.Context, args []string) error {
            pageName := args[0]
            // Validate that the page exists
            if !pageExists(ctx.Pages, pageName) {
                return fmt.Errorf("page %q does not exist", pageName)
            }
            return nil
        },
        func(ctx *template.Context, pageName string) {
            ctx.Pages.SwitchToPage(pageName)
        },
    ).
    Build()
```

**Note**: The validator is only called after argument count validation passes. Use it for semantic validation like checking if a page exists or validating argument format.

## Package Structure

```
github.com/cassdeckard/tviewyaml/
├── app.go           # Application builder (deprecated CreateApp)
├── builder.go       # AppBuilder with Builder pattern API
├── builder/         # UI builder components
│   ├── builder.go
│   ├── callbacks.go
│   ├── factory.go
│   └── properties.go
├── config/          # Configuration loading and types
│   ├── loader.go
│   ├── types.go
│   └── validator.go
├── keys/            # Key binding parsing
│   └── keys.go      # ParseKey for key string parsing
└── template/        # Template execution
    ├── builtins.go  # Built-in template functions
    ├── context.go
    ├── executor.go
    ├── keybinding.go  # MatchesKeyBinding for key event matching
    └── registry.go  # Function registry system
```

## Requirements

- Go 1.21 or higher
- [tview](https://github.com/rivo/tview) - Terminal UI library
- [tcell](https://github.com/gdamore/tcell) - Terminal cell library
- [yaml.v3](https://gopkg.in/yaml.v3) - YAML parser

## License

[MIT License](LICENSE)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Publishing

This package is published at `github.com/cassdeckard/tviewyaml`.

To use it in your projects:
```bash
go get github.com/cassdeckard/tviewyaml
```

### For Package Maintainers

To publish a new version:
1. Tag the release: `git tag v0.1.0`
2. Push the tag: `git push origin v0.1.0`
3. Go will automatically make it available via the module proxy

## Credits

Built on top of the excellent [tview](https://github.com/rivo/tview) library by rivo.
