# Contributing to Scribe

Thank you for your interest in contributing!

## Development Setup

1. Clone the repository:
```bash
git clone https://github.com/yusnelgg/scribe.git
cd scribe
```

2. Install dependencies:
```bash
go mod download
```

3. Run tests:
```bash
go test ./...
```

4. Build:
```bash
go build -o scribe ./cmd/scribe
```

## Code Style

- Follow Go idioms and `gofmt`
- Add tests for new features
- Keep functions small and focused
- Document exported functions

## Commit Messages

Use conventional commits:

```
feat(parser): add support for router groups
fix(scanner): ignore vendor directories
docs: update README with new examples
test: add test for Gin route detection
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`

## Pull Request Process

1. Fork and create a branch from `main`
2. Make your changes with tests
3. Ensure all tests pass
4. Update documentation if needed
5. Submit a clear PR description

## Adding Framework Support

To add support for a new framework:

1. Create `internal/parser/<framework>.go`
2. Implement the `Parser` interface:
```go
type Parser interface {
    Parse(file string) ([]Route, error)
}
```

3. Add framework detection in `cmd/scribe/main.go`
4. Add tests
5. Update README.md

## Questions?

Open an issue for bugs or feature requests.
