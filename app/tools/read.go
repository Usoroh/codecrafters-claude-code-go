package tools

import (
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

type ReadTool struct{}

// GetTool returns the tool definition for reading file contents.
func (rt ReadTool) GetTool() openai.ChatCompletionToolUnionParam {
	return openai.ChatCompletionToolUnionParam{
		OfFunction: &openai.ChatCompletionFunctionToolParam{
			Function: shared.FunctionDefinitionParam{
				Name:        Read,
				Description: openai.String("Read and return the contents of a file"),
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"file_path": map[string]string{
							"type":        "string",
							"description": "The path to the file to read",
						},
					},
					"required": []string{"file_path"},
				},
			},
		},
	}
}

// Execute reads the file provided in arguments
func (rt ReadTool) Execute(args map[string]any) (string, error) {
	fileName, err := getStringArg(args, "file_path")
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("file_path [%s]: %w", fileName, err)
	}

	return string(content), nil
}
