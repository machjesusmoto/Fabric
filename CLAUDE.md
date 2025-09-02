# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

Fabric is an open-source AI framework written in Go that helps users augment their capabilities using AI through organized prompts called "Patterns". It integrates with multiple AI providers (OpenAI, Anthropic, Ollama, Google Gemini, AWS Bedrock, Perplexity) and provides both CLI and web interfaces.

## Common Development Commands

### Building and Running

```bash
# Install from source
go install github.com/danielmiessler/fabric/cmd/fabric@latest

# Run setup (required after first installation)
fabric --setup

# Basic usage
fabric --pattern summarize < input.txt
pbpaste | fabric --pattern extract_wisdom --stream

# Run with specific model
fabric --pattern analyze_claims --model gpt-4 --stream

# YouTube transcript processing
fabric -y "https://youtube.com/watch?v=VIDEO_ID" --pattern extract_wisdom

# Web scraping with Jina AI
fabric -u https://example.com -p analyze_claims

# Run the web server
fabric --serve  # Default API server on :8080
fabric --serveOllama  # With Ollama endpoints
```

### Helper Tools

```bash
# Install helper tools
go install github.com/danielmiessler/fabric/cmd/code_helper@latest
go install github.com/danielmiessler/fabric/cmd/to_pdf@latest

# Use code_helper for AI-assisted coding
code_helper . "Add error handling to all functions" | fabric --pattern create_coding_feature

# Convert LaTeX to PDF
echo "AI security primer" | fabric --pattern write_latex | to_pdf
```

### Web Interface Development

```bash
cd web
pnpm install  # or npm install
pnpm dev      # Start development server
pnpm build    # Build for production
pnpm check    # Type checking
pnpm lint     # Lint code
pnpm format   # Format code
pnpm test     # Run tests
```

### Testing

```bash
# Run all Go tests
go test -v ./...

# Run specific package tests
go test -v ./internal/cli/...
go test -v ./internal/plugins/ai/...

# Check formatting with Nix
nix flake check
```

## High-Level Architecture

### Core Components

1. **CLI Interface** (`internal/cli/`)
   - Main entry point and flag parsing
   - Command routing and execution
   - Setup wizard and configuration management

2. **Plugin System** (`internal/plugins/`)
   - **AI Providers** (`ai/`): Modular vendor implementations (OpenAI, Anthropic, Ollama, etc.)
   - **Database** (`db/`): File-based pattern and context storage
   - **Strategies** (`strategy/`): Prompt modification strategies (CoT, ToT, reflexion)
   - **Templates** (`template/`): Extension system for custom functionality

3. **Core Domain** (`internal/core/`)
   - Chatter interface for AI interactions
   - Plugin registry for dynamic vendor loading
   - Message handling and streaming

4. **Tools** (`internal/tools/`)
   - YouTube transcript extraction
   - HTML readability conversion
   - Jina AI web scraping integration
   - Pattern loading and management

5. **Data** (`data/`)
   - **Patterns**: 200+ reusable AI prompts organized by task
   - **Strategies**: JSON-based prompt modification strategies

### Key Design Patterns

- **Plugin Architecture**: All AI vendors implement the `core.Chatter` interface
- **Registry Pattern**: Dynamic vendor registration and selection
- **Streaming Support**: Real-time response streaming for better UX
- **Extension System**: Custom templates and functions via YAML configuration
- **Pattern Organization**: Markdown-based prompts with system/user sections

### AI Provider Integration

Each AI provider (`internal/plugins/ai/`) implements:
- `NewClient()` - Initialize with configuration
- `Send()` - Process messages and return responses
- `SendStream()` - Stream responses in real-time
- Model listing and validation
- Provider-specific features (image generation, TTS, web search)

### Web Interface Architecture

The web interface (SvelteKit + TypeScript) provides:
- Pattern browsing and execution
- Session management
- PDF conversion from various formats
- YouTube transcript integration
- Obsidian.md vault integration for note-taking

## Configuration

Fabric uses environment variables and config files:
- `~/.config/fabric/` - Main configuration directory
- `~/.config/fabric/patterns/` - User patterns
- `~/.config/fabric/.env` - API keys and settings
- Custom patterns directory (configurable via setup)

## Key Files to Understand

- `cmd/fabric/main.go` - CLI entry point
- `internal/cli/cli.go` - Main CLI logic and command routing
- `internal/core/chatter.go` - AI provider interface definition
- `internal/plugins/ai/vendor.go` - Vendor management
- `internal/server/serve.go` - REST API server
- `web/src/routes/+page.svelte` - Web UI main page

## Pattern Development

Patterns are markdown files in `data/patterns/[pattern_name]/system.md`:
- Use clear, structured markdown
- Focus on system prompts for consistency
- Include step-by-step instructions
- Use formatting to emphasize important points

## Extension System

Create custom extensions via YAML configuration:
- Define custom template functions
- Add shell command integrations
- Create reusable templates
- Support for SQLite queries and external APIs