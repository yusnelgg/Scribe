# Scribe

> Generate API documentation, tests, and analysis reports from your codebase.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)

Scribe is a CLI tool that analyzes your backend codebase and automatically generates:

- **OpenAPI 3.0 documentation** - Standard API specification
- **Test files** - Ready-to-use test templates  
- **Analysis reports** - Route inventory and metrics

## Features

- Zero configuration required
- Static code analysis (no AI required by default)
- Pluggable AI enhancement (optional Ollama/OpenAI support)
- Focus on Gin framework (extensible architecture)
- Clean, idiomatic Go code

## Installation

### From Source

```bash
git clone https://github.com/yusnel/mt-scribe.git
cd scribe
go install
```

### Pre-built Binaries

Download from [Releases](https://github.com/yusnel/mt-scribe/releases) for your platform.

## Quick Start

```bash
# Scan current directory
scribe scan .

# Scan specific project
scribe scan ./my-api

# With verbose output
scribe scan . -v

# Enable AI enhancement (requires Ollama)
scribe scan . -ai
```

## Generated Files

After scanning, you'll find:

```
.
├── openapi.json    # OpenAPI 3.0 specification
├── report.md       # Analysis report
└── tests/          # Test templates
    ├── test_get_users.go
    ├── test_post_users.go
    └── ...
```

## CLI Options

| Flag | Description | Default |
|------|-------------|---------|
| `-path` | Project path to scan | `.` |
| `-framework` | Web framework (gin) | `gin` |
| `-ai` | Enable AI enhancement | `false` |
| `-ai-provider` | AI provider (ollama) | `ollama` |
| `-ai-endpoint` | AI endpoint URL | `http://localhost:11434` |
| `-output` | Output directory | `.` |
| `-v` | Verbose output | `false` |

## Architecture

```
scribe/
├── cmd/scribe/          # CLI entry point
└── internal/
    ├── scanner/         # File system scanning
    ├── parser/          # Code analysis (AST)
    ├── generator/       # Output generation
    ├── ai/              # AI integration
    └── formatter/       # Output formatting
```

## Supported Frameworks

- [x] Gin (Go)
- [ ] Echo (Go)
- [ ] Fiber (Go)
- [ ] Chi (Go)
- [ ] Express (Node.js)
- [ ] FastAPI (Python)

## AI Integration

Scribe supports optional AI enhancement. By default, it works purely through static analysis.

When AI is enabled, Scribe:
- **Enhances tests** - Generates comprehensive tests with proper setup, mocks, and assertions
- **Improves OpenAPI specs** - Adds descriptions, summaries, and request/response schemas

### Ollama (Local, Recommended)

```bash
# Install Ollama
brew install ollama  # macOS
curl -fsSL https://ollama.com/install.sh | sh  # Linux

# Start Ollama
ollama serve

# Use with scribe
scribe scan . -ai
```

### OpenAI

```bash
export OPENAI_API_KEY=your-key
scribe scan . -ai -ai-provider=openai
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md).

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE)

## Roadmap

- [ ] Extract request/response types from function signatures
- [ ] Detect middleware usage
- [ ] Support router groups
- [ ] Generate Postman collections
- [ ] Webhook detection
- [ ] Authentication/authorization analysis
