package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/codecrafters-io/claude-code-starter-go/app/tools"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	var prompt string
	flag.StringVar(&prompt, "p", "", "Prompt to send to LLM")
	flag.Parse()

	if prompt == "" {
		panic("Prompt must not be empty")
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	baseUrl := os.Getenv("OPENROUTER_BASE_URL")
	if baseUrl == "" {
		baseUrl = "https://openrouter.ai/api/v1"
	}

	if apiKey == "" {
		panic("Env variable OPENROUTER_API_KEY not found")
	}

	readTool := tools.ReadTool{}

	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseUrl))
	resp, err := client.Chat.Completions.New(context.Background(),
		openai.ChatCompletionNewParams{
			Model: "anthropic/claude-haiku-4.5",
			Messages: []openai.ChatCompletionMessageParamUnion{
				{
					OfUser: &openai.ChatCompletionUserMessageParam{
						Content: openai.ChatCompletionUserMessageParamContentUnion{
							OfString: openai.String(prompt),
						},
					},
				},
			},
			Tools: []openai.ChatCompletionToolUnionParam{readTool.GetTool()},
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if len(resp.Choices) == 0 {
		panic("No choices in response")
	}

	// handle tool calls
	if len(resp.Choices[0].Message.ToolCalls) > 0 {
		name := resp.Choices[0].Message.ToolCalls[0].Function.Name
		if name == tools.Read {
			rawArgs := resp.Choices[0].Message.ToolCalls[0].Function.Arguments
			var args map[string]any

			err := json.Unmarshal([]byte(rawArgs), &args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error unmarshalling tool arguments: %v\n", err)
			}

			result, err := readTool.Execute(args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing tool: %v\n", err)
			}

			fmt.Println(result)
		}
	}

	// TODO: Uncomment the line below to pass the first stage
	fmt.Print(resp.Choices[0].Message.Content)
}
