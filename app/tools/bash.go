package tools

import (
	"fmt"
	"os/exec"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/shared"
)

type BashTool struct{}

func init() { Register("Bash", BashTool{}) }

// GetTool returns the tool definition for executing shell commands.
func (bt BashTool) GetTool() openai.ChatCompletionToolUnionParam {
	return openai.ChatCompletionToolUnionParam{
		OfFunction: &openai.ChatCompletionFunctionToolParam{
			Function: shared.FunctionDefinitionParam{
				Name:        "Bash",
				Description: openai.String("Execute a shell command"),
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"command": map[string]string{
							"type":        "string",
							"description": "The command to execute",
						},
					},
					"required": []string{"command"},
				},
			},
		},
	}
}

// Execute executes the command using shell
func (bt BashTool) Execute(args map[string]any) (string, error) {
	command, err := getStringArg(args, "command")
	if err != nil {
		return "", err
	}

	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command [%s]: %s\n%s", command, err, string(output))
	}

	return string(output), nil
}
