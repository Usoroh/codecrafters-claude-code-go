package tools

import (
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
)

// Tool defines the interface that all tools must implement.
type Tool interface {
	GetTool() openai.ChatCompletionToolUnionParam
	Execute(args map[string]any) (string, error)
}

var registry = map[string]Tool{}

// Register adds a tool to the global registry
func Register(name string, t Tool) {
	registry[name] = t
}

// AllTools returns the tool schemas for all registered tools.
func AllTools() []openai.ChatCompletionToolUnionParam {
	params := make([]openai.ChatCompletionToolUnionParam, 0, len(registry))
	for _, t := range registry {
		params = append(params, t.GetTool())
	}

	return params
}

// GetToolCallResult executes a tool call and returns its result.
func GetToolCallResult(toolCall openai.ChatCompletionMessageToolCallUnion) (string, error) {
	var args map[string]any
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		return "", fmt.Errorf("error unmarshalling tool arguments: %w", err)
	}

	tool, ok := registry[toolCall.Function.Name]
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
