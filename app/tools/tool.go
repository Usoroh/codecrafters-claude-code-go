package tools

import (
	"encoding/json"
	"fmt"
	"github.com/openai/openai-go/v3"
)

const (
	Read  = "Read"
	Write = "Write"
	Bash  = "Bash"
)

type Tool interface {
	GetTool() openai.ChatCompletionToolUnionParam
	Execute(args map[string]any) (string, error)
}

// GetToolCallResult return the result of a tool call
func GetToolCallResult(toolCall openai.ChatCompletionMessageToolCallUnion) (string, error) {
	var args map[string]any
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		return "", fmt.Errorf("error unmarshalling tool arguments: %w", err)
	}

	tool, ok := getTool(toolCall.Function.Name)
	if !ok {
		return "", fmt.Errorf("unknown tool: %s", toolCall.Function.Name)
	}

	result, err := tool.Execute(args)
	if err != nil {
		return "", fmt.Errorf("error executing tool: %w", err)
	}

	return result, nil
}

// getStringArg extracts a required string argument from the args map.
func getStringArg(args map[string]any, key string) (string, error) {
	val, ok := args[key]
	if !ok {
		return "", fmt.Errorf("missing required argument: %s", key)
	}

	s, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("%s must be a string", key)
	}

	return s, nil
}

// getTool returns a tool when given its name
func getTool(name string) (Tool, bool) {
	toolMap := map[string]Tool{
		Read:  ReadTool{},
		Write: WriteTool{},
		Bash:  BashTool{},
	}

	tool, ok := toolMap[name]

	return tool, ok
}
