package tools

import (
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

type WriteTool struct{}

// GetTool returns the tool definition for reading file contents.
func (wt WriteTool) GetTool() openai.ChatCompletionToolUnionParam {
	return openai.ChatCompletionToolUnionParam{
		OfFunction: &openai.ChatCompletionFunctionToolParam{
			Function: shared.FunctionDefinitionParam{
				Name:        Write,
				Description: openai.String("Write content to a file"),
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"file_path": map[string]string{
							"type":        "string",
							"description": "The path of the file to write to",
						},
						"content": map[string]string{
							"type":        "string",
							"description": "The content to write to the file",
						},
					},
					"required": []string{"file_path", "content"},
				},
			},
		},
	}
}

// Execute writes content to a file
func (wt WriteTool) Execute(args map[string]any) (string, error) {
	fileName, err := getStringArg(args, "file_path")
	if err != nil {
		return "", err
	}

	content, err := getStringArg(args, "content")
	if err != nil {
		return "", err
	}

	err = os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		return "", fmt.Errorf("file_path [%s]: %w", fileName, err)
	}

	return fmt.Sprintf("Written following content to %s:%s", fileName, content), nil
}
