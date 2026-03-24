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

	messages := []openai.ChatCompletionMessageParamUnion{
		{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: openai.String(prompt),
				},
			},
		},
	}

	for {
		resp, err := client.Chat.Completions.New(context.Background(),
			openai.ChatCompletionNewParams{
				Model:    "anthropic/claude-haiku-4.5",
				Messages: messages,
				Tools:    []openai.ChatCompletionToolUnionParam{readTool.GetTool()},
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		if len(resp.Choices) == 0 {
			panic("No choices in response")
		}

		message := resp.Choices[0].Message

		messages = append(messages, message.ToParam())

		if len(message.ToolCalls) == 0 {
			fmt.Fprintf(os.Stdout, message.Content)
			os.Exit(0)
		}

		for _, toolCall := range message.ToolCalls {
			content, err := getToolCallResult(toolCall)
			if err != nil {
				content = err.Error()
			}

			messages = append(messages, openai.ChatCompletionMessageParamUnion{
				OfTool: &openai.ChatCompletionToolMessageParam{
					Content: openai.ChatCompletionToolMessageParamContentUnion{
						OfString: openai.String(content),
					},
				},
			})
		}
	}
}

// getToolCallResult return the result of a tool call
func getToolCallResult(toolCall openai.ChatCompletionMessageToolCallUnion) (string, error) {
	var args map[string]any
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		return "", fmt.Errorf("error unmarshalling tool arguments: %v\n", err)
	}

	tool, ok := getTool(toolCall.Function.Name)
	if !ok {
		return "", fmt.Errorf("error getting tool: %v\n", toolCall.Function.Name)
	}

	result, err := tool.Execute(args)
	if err != nil {
		return "", fmt.Errorf("error executing tool: %v\n", err)
	}

	fmt.Println("RESULT: ", result)

	return result, nil
}

// getTool returns a tool when given its name
func getTool(name string) (tools.Tool, bool) {
	toolMap := map[string]tools.Tool{
		tools.Read: tools.ReadTool{},
	}

	tool, ok := toolMap[name]
	return tool, ok
}
