# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a CodeCrafters "Build Your Own Claude Code" challenge solution in Go. The program implements an LLM-powered coding assistant that communicates with an OpenAI-compatible API (OpenRouter) and supports tool calling.

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

- **`app/main.go`**: Entry point. Parses `-p` flag for the prompt, creates an OpenAI client pointed at OpenRouter, sends chat completion requests with tool definitions, and prints the response.
- **`app/tools/schema.go`**: Defines the JSON-serializable types (`Tool`, `Function`, `Parameters`, `Property`) for OpenAI-compatible tool/function calling schema.
- **`app/tools/read.go`**: Defines the `Read` tool (file reading capability) using the schema types.

The program uses `github.com/openai/openai-go/v3` to talk to the OpenRouter API with model `anthropic/claude-haiku-4.5`. Tools are registered by passing them in the `Tools` field of the chat completion request. New tools should follow the pattern in `read.go`: define a `GetXTool()` function returning a `Tool` struct.
