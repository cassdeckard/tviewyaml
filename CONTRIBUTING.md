# Contributing to tviewyaml

Thank you for your interest in contributing to tviewyaml!

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/cassdeckard/tviewyaml.git
   cd tviewyaml
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   go build ./...
   ```

4. Run tests (if available):
   ```bash
   go test ./...
   ```

## Testing Your Changes

You can test your changes using the example application:

```bash
cd example
go run main.go
```

Or create your own test YAML configurations.

## Code Style

- Follow standard Go formatting (`gofmt`)
- Add comments for exported functions and types
- Keep functions focused and single-purpose
- Use meaningful variable names

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Commit your changes (`git commit -m 'Add amazing feature'`)
5. Push to your branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## Adding New Features

When adding new widget types or features:

1. Update `config/types.go` with any new configuration fields
2. Update `builder/factory.go` if adding new widget types
3. Update `builder/properties.go` to handle new properties
4. Update `README.md` with documentation
5. Add examples in the `example/` directory if applicable

## Reporting Bugs

Please open an issue with:
- A clear description of the bug
- Steps to reproduce
- Expected behavior
- Actual behavior
- Your environment (Go version, OS, etc.)

## Feature Requests

Feature requests are welcome! Please open an issue describing:
- The feature you'd like to see
- Why it would be useful
- Any implementation ideas you have

## Questions?

Feel free to open an issue for any questions about contributing!
