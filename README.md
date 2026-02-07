# tviewyaml

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

### 1. Create a Root Configuration (`root.yaml`)

```yaml
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
    // Create app from YAML config directory
    app, err := tviewyaml.CreateApp("./config")
    if err != nil {
        log.Fatal(err)
    }

    // Run the application
    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## Configuration Examples

### Flex Layout with Multiple Widgets

```yaml
type: flex
direction: row
border: true
title: Dashboard
items:
  - primitive:
      type: textView
      text: "Welcome to the app!"
      border: true
      title: Info
    fixedSize: 0
    proportion: 1
    focus: false
  
  - primitive:
      type: list
      border: true
      title: Actions
      listItems:
        - mainText: Action 1
          onSelected: "{{ showSimpleModal \"Action 1 clicked\" }}"
        - mainText: Action 2
          onSelected: "{{ showSimpleModal \"Action 2 clicked\" }}"
    fixedSize: 0
    proportion: 1
    focus: true
```

### Form Example

```yaml
type: form
title: User Input
border: true
formItems:
  - type: inputfield
    label: "Name:"
    value: ""
    fieldWidth: 20
  
  - type: inputfield
    label: "Age:"
    value: ""
    fieldWidth: 3
    acceptanceFunc: integer
  
  - type: checkbox
    label: "Subscribe to newsletter"
    checked: false
  
  - type: dropdown
    label: "Country:"
    options:
      - USA
      - Canada
      - UK
      - Other
  
  - type: button
    label: Submit
    onSelected: "{{ showSimpleModal \"Form submitted!\" \"OK\" }}"
  
  - type: button
    label: Cancel
    onSelected: "{{ switchToPage \"main\" }}"
```

### Table Example

```yaml
type: table
title: Data Table
border: true
tableData:
  headers:
    - Name
    - Age
    - City
  rows:
    - ["Alice", "30", "New York"]
    - ["Bob", "25", "San Francisco"]
    - ["Charlie", "35", "Chicago"]
```

## Template Functions

Built-in template functions for callbacks:

- `switchToPage "pageName"` - Navigate to a different page
- `removePage "pageName"` - Remove a page from the stack
- `stopApp` - Exit the application
- `showSimpleModal "text" "button1" "button2"` - Show a modal dialog
- `noop` - No operation (placeholder callback)

## Package Structure

```
github.com/cassdeckard/tviewyaml/
├── app.go           # Main application entry point
├── builder/         # UI builder components
│   ├── builder.go
│   ├── callbacks.go
│   ├── factory.go
│   └── properties.go
├── config/          # Configuration loading and types
│   ├── loader.go
│   ├── types.go
│   └── validator.go
└── template/        # Template execution
    ├── context.go
    ├── executor.go
    └── functions.go
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
