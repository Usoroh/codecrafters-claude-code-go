# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a CodeCrafters "Build Your Own Claude Code" challenge solution in Go. The program implements an LLM-powered coding assistant that communicates with an OpenAI-compatible API (OpenRouter) and supports tool calling (Read, Write, Bash).

## Build & Run Commands

```sh
# Build and run locally
./your_program.sh -p "your prompt here"

# Build only
go build -o /tmp/codecrafters-build-claude-code-go app/*.go

# Submit to CodeCrafters
codecrafters submit
```

## Environment Variables

- `OPENROUTER_API_KEY` (required): API key for OpenRouter
- `OPENROUTER_BASE_URL` (optional): defaults to `https://openrouter.ai/api/v1`

## Architecture

- **`app/main.go`**: Entry point. Parses `-p` flag, creates an OpenAI client pointed at OpenRouter, and runs the agentic loop: sends chat completions, executes any tool calls, appends results to the conversation, and repeats until the model returns a plain text response.
- **`app/tools/tool.go`**: Defines the `Tool` interface (`GetTool()` + `Execute()`), a global name→tool registry, and the `GetToolCallResult` dispatcher that unmarshals arguments and routes to the correct tool. Also has `getStringArg` helper for extracting required string args.
- **`app/tools/{read,write,bash}.go`**: Individual tool implementations. Each is a struct implementing the `Tool` interface, self-registering via `init()` + `Register()`.

### Adding a new tool

1. Create `app/tools/mytool.go` with a struct implementing `Tool` (methods: `GetTool()`, `Execute(args map[string]any)`).
2. Add `func init() { Register("MyTool", MyToolStruct{}) }` — this auto-registers it at startup.
3. No changes needed in `main.go`; `tools.AllTools()` picks up all registered tools.

The program uses `github.com/openai/openai-go/v3` to talk to the OpenRouter API with model `anthropic/claude-haiku-4.5`.
