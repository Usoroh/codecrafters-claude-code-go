# Build Your Own Claude Code (Go)

A minimal LLM-powered coding assistant built in Go, inspired by [Claude Code](https://claude.ai/code). This project is a solution to the [CodeCrafters "Build Your Own Claude Code" challenge](https://codecrafters.io/challenges/claude-code).

The assistant sends prompts to an LLM via an OpenAI-compatible API (OpenRouter), supports tool calling, and runs in an agentic loop — allowing the model to read files, write files, and execute shell commands autonomously until it produces a final response.

## Tools

| Tool | Description |
|------|-------------|
| **Read** | Reads and returns the contents of a file |
| **Write** | Writes content to a file |
| **Bash** | Executes a shell command and returns the output |

## Prerequisites

- Go 1.26+
- An [OpenRouter](https://openrouter.ai/) API key

## Setup

1. Clone the repository:

   ```sh
   git clone https://github.com/Usoroh/codecrafters-claude-code-go.git
   cd codecrafters-claude-code-go
   ```

2. Set your API key:

   ```sh
   export OPENROUTER_API_KEY="your-api-key-here"
   ```

3. (This is optional) Override the base URL if using a different OpenAI-compatible provider:

   ```sh
   export OPENROUTER_BASE_URL="https://your-provider.com/api/v1"
   ```

## Usage

```sh
./your_program.sh -p "your prompt here"
```

### Examples

```sh
# Ask it to read a file
./your_program.sh -p "What does main.go do?"

# Ask it to create a file
./your_program.sh -p "Create a hello.txt file that says hello world"

# Ask it to run a command
./your_program.sh -p "List all Go files in this project"
```

## Building Manually

```sh
go build -o /tmp/codecrafters-build-claude-code-go app/*.go
/tmp/codecrafters-build-claude-code-go -p "your prompt"
```

## Project Structure

```
app/
├── main.go            # Entry point: parses flags, runs the agent loop
└── tools/
    ├── tool.go        # Tool interface, dispatch, and argument helpers
    ├── read.go        # Read tool implementation
    ├── write.go       # Write tool implementation
    └── bash.go        # Bash tool implementation
```

## How It Works

1. The user provides a prompt via the `-p` flag.
2. The prompt is sent to the LLM (`anthropic/claude-haiku-4.5` on OpenRouter).
3. If the model responds with tool calls, each tool is executed and the results are appended to the conversation.
4. The loop continues until the model produces a plain text response, which is printed to stdout.
