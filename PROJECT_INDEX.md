# Fabric AI Framework - Project Index

## Overview
Fabric is an open-source AI augmentation framework written in Go that provides organized prompts (called "Patterns") and integrations with multiple AI providers. It addresses the AI integration problem by creating a unified interface for AI capabilities.

**Version**: v1.4.272  
**Language**: Go 1.19+  
**License**: MIT  
**Mission**: Human flourishing via AI augmentation

## Quick Links
- [Installation](#installation-setup)
- [Core Components](#core-components)
- [API Reference](#api-reference)
- [Patterns Library](#patterns-library)
- [Configuration](#configuration)
- [Development](#development)

## Architecture Overview

```
fabric/
├── cmd/                    # Command-line tools and entry points
├── internal/               # Core implementation
│   ├── cli/               # CLI interface and flag handling
│   ├── core/              # Core domain interfaces
│   ├── domain/            # Domain models and business logic
│   ├── plugins/           # Plugin system (AI, DB, strategies)
│   ├── server/            # REST API server
│   └── tools/             # Utility tools
├── data/                  # Patterns and strategies
│   ├── patterns/          # 200+ AI prompt patterns
│   └── strategies/        # Reasoning strategies
└── web/                   # Web interface (SvelteKit)
```

## Core Components

### 1. Command-Line Interface (`internal/cli/`)

#### Primary Files
- **`cli.go`** - Main CLI orchestration and command routing
- **`flags.go`** - Flag parsing and configuration management
- **`chat.go`** - Chat processing and AI interactions
- **`initialization.go`** - Fabric initialization and setup

#### Key Functions
| Function | File | Purpose |
|----------|------|---------|
| `Cli()` | cli.go:16 | Main entry point for CLI |
| `Init()` | flags.go:103 | Initialize flags and configuration |
| `handleChatProcessing()` | chat.go:15 | Process chat interactions |
| `initializeFabric()` | initialization.go:16 | Initialize core framework |

#### Command Categories
- **Setup & Configuration**: `--setup`, `--setupOllama`, `--changeDefaultModel`
- **Pattern Management**: `--listPatterns`, `--listAllPatterns`, `--updatePatterns`
- **Chat Operations**: `--pattern`, `--model`, `--stream`, `--temperature`
- **Output Control**: `--output`, `--copy`, `--markdown`, `--json`
- **Server Mode**: `--serve`, `--serveOllama`

### 2. Plugin System (`internal/plugins/`)

#### AI Providers (`ai/`)
Modular implementations for different AI vendors:

| Provider | File | Features |
|----------|------|----------|
| **OpenAI** | `openai/openai.go` | GPT models, DALL-E, TTS |
| **Anthropic** | `anthropic/anthropic.go` | Claude models, OAuth support |
| **Google Gemini** | `gemini/gemini.go` | Gemini models, voices |
| **Ollama** | `ollama/ollama.go` | Local model support |
| **AWS Bedrock** | `bedrock/bedrock.go` | AWS AI services |
| **Azure** | `azure/azure.go` | Azure OpenAI integration |
| **Perplexity** | `perplexity/perplexity.go` | Web-aware AI |

#### Database (`db/`)
File-based storage system:
- **`fsdb/db.go`** - Main database interface
- **`fsdb/patterns.go`** - Pattern storage and retrieval
- **`fsdb/contexts.go`** - Context management
- **`fsdb/sessions.go`** - Session handling

#### Strategies (`strategy/`)
Prompt modification strategies for enhanced reasoning:
- Chain of Thought (CoT)
- Tree of Thought (ToT)
- Reflexion
- Self-consistency

### 3. Core Domain (`internal/core/`)

#### Key Interfaces
```go
// Chatter - Core AI provider interface
type Chatter interface {
    Send(ctx context.Context, req *ChatRequest, cfg *ChatOptions) (*ChatResponse, error)
    SendStream(ctx context.Context, req *ChatRequest, cfg *ChatOptions) (<-chan ChatStreamResponse, error)
    GetModels() []string
}
```

#### Plugin Registry
- **`plugin_registry.go`** - Dynamic vendor registration
- **`chatter.go`** - AI provider interface definition

### 4. Server API (`internal/server/`)

REST API endpoints for web integration:

| Endpoint | Handler | Purpose |
|----------|---------|---------|
| `/patterns` | `patterns.go` | Pattern management |
| `/models` | `models.go` | Model listing |
| `/chat` | `chat.go` | Chat interactions |
| `/contexts` | `contexts.go` | Context handling |
| `/youtube` | `youtube.go` | YouTube transcripts |

### 5. Tools (`internal/tools/`)

Utility modules for enhanced functionality:

- **YouTube Integration** (`youtube/`) - Transcript extraction
- **HTML Conversion** (`converter/`) - Readability processing
- **Jina AI** (`jina/`) - Web scraping
- **Pattern Loader** (`patterns_loader.go`) - Dynamic pattern loading
- **Custom Patterns** (`custom_patterns/`) - User pattern management

## Patterns Library

### Pattern Structure
Each pattern is a markdown file with structured prompts:
```
data/patterns/[pattern_name]/
├── system.md    # System prompt (required)
├── user.md      # User prompt template (optional)
└── README.md    # Pattern documentation (optional)
```

### Pattern Categories (200+ patterns)

#### Analysis Patterns
- `analyze_claims` - Verify claims and statements
- `analyze_paper` - Academic paper analysis
- `analyze_risk` - Risk assessment
- `analyze_threat_report` - Security threat analysis
- `analyze_code` - Code review and analysis

#### Creation Patterns
- `create_coding_project` - Generate project structures
- `create_git_diff_commit` - Create commit messages
- `create_documentation` - Technical documentation
- `create_api` - API endpoint design
- `create_test_suite` - Test generation

#### Extraction Patterns
- `extract_wisdom` - Key insights extraction
- `extract_article_wisdom` - Article summarization
- `extract_predictions` - Future predictions
- `extract_requirements` - Requirements gathering

#### Improvement Patterns
- `improve_writing` - Writing enhancement
- `improve_code` - Code optimization
- `improve_prompt` - Prompt refinement

### Strategy System

Reasoning strategies modify prompts for better results:
```json
{
  "name": "Chain of Thought",
  "prefix": "Let's think step by step...",
  "suffix": "Therefore, the answer is..."
}
```

## API Reference

### CLI Commands

```bash
# Basic usage
fabric --pattern summarize < input.txt
fabric --pattern extract_wisdom --stream

# YouTube processing
fabric -y "https://youtube.com/watch?v=VIDEO_ID" --pattern summarize

# Web scraping
fabric -u https://example.com --pattern analyze_claims

# Model selection
fabric --pattern create_code --model gpt-4 --temperature 0.7

# Server mode
fabric --serve               # API server on :8080
fabric --serveOllama         # With Ollama endpoints
```

### Go Package Usage

```go
import (
    "github.com/danielmiessler/fabric/internal/core"
    "github.com/danielmiessler/fabric/internal/cli"
)

// Initialize Fabric
registry, err := core.NewPluginRegistry()
flags := cli.Init()

// Create chat request
request := &domain.ChatRequest{
    Messages: []Message{
        {Role: "user", Content: "Hello"},
    },
}

// Send to AI provider
response, err := chatter.Send(ctx, request, options)
```

## Configuration

### Environment Variables
Location: `~/.config/fabric/.env`
```bash
# API Keys
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
GOOGLE_API_KEY=...

# Model Configuration
DEFAULT_MODEL=gpt-4
DEFAULT_TEMPERATURE=0.7

# Paths
PATTERNS_PATH=~/.config/fabric/patterns
CUSTOM_PATTERNS_PATH=~/my-patterns
```

### Directory Structure
```
~/.config/fabric/
├── .env                    # API keys and settings
├── patterns/               # Local patterns
├── contexts/              # Saved contexts
└── sessions/              # Session history
```

## Development

### Building from Source

```bash
# Install from source
go install github.com/danielmiessler/fabric/cmd/fabric@latest

# Build locally
go build -o fabric cmd/fabric/main.go

# Run tests
go test -v ./...

# Format code
go fmt ./...
```

### Web Interface Development

```bash
cd web
pnpm install    # Install dependencies
pnpm dev        # Development server
pnpm build      # Production build
pnpm test       # Run tests
```

### Testing

#### Go Tests
```bash
# All tests
go test -v ./...

# Specific package
go test -v ./internal/cli/...

# With coverage
go test -cover ./...
```

#### Pattern Testing
```bash
# Test pattern execution
echo "Test input" | fabric --pattern summarize --dryRun

# Validate pattern syntax
fabric --validatePattern ./data/patterns/my_pattern/
```

## Extension System

### Custom Template Functions

Create YAML configuration for custom functionality:
```yaml
name: my_extension
functions:
  - name: custom_func
    command: "echo 'Custom output'"
templates:
  - name: my_template
    content: "{{.Input | custom_func}}"
```

### Plugin Development

Implement the `Chatter` interface:
```go
type MyProvider struct{}

func (m *MyProvider) Send(ctx context.Context, req *ChatRequest, cfg *ChatOptions) (*ChatResponse, error) {
    // Implementation
}

func (m *MyProvider) GetModels() []string {
    return []string{"my-model-1", "my-model-2"}
}
```

## Helper Tools

### code_helper
AI-assisted coding tool:
```bash
code_helper . "Add error handling" | fabric --pattern create_coding_feature
```

### to_pdf
Convert outputs to PDF:
```bash
echo "Document title" | fabric --pattern write_latex | to_pdf
```

### generate_changelog
Automated changelog generation:
```bash
cd cmd/generate_changelog
go run main.go --since v1.4.0
```

## Project Statistics

- **Go Files**: 100+ source files
- **Patterns**: 200+ AI prompts
- **Strategies**: 9 reasoning strategies
- **AI Providers**: 10+ integrations
- **Lines of Code**: ~15,000+ Go, ~5,000+ TypeScript
- **Test Coverage**: Comprehensive unit and integration tests

## Navigation Map

### Quick Access Points
- [Main Entry](cmd/fabric/main.go) - CLI entry point
- [CLI Logic](internal/cli/cli.go) - Core CLI implementation
- [AI Interface](internal/core/chatter.go) - Provider interface
- [Pattern Directory](data/patterns/) - All patterns
- [Web UI](web/src/routes/+page.svelte) - Main web interface

### Module Dependencies
```
cmd/fabric → internal/cli → internal/core
                          ↓
                    internal/plugins/ai
                          ↓
                    internal/domain
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## License

MIT License - See [LICENSE](LICENSE) for details.

---
*Generated with Fabric v1.4.272 - Last updated: 2025*