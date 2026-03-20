package tools

import "github.com/openai/openai-go/v3"

const Read = "Read"

type Tool interface {
	GetTool() openai.ChatCompletionToolUnionParam
	Execute(args map[string]any) (string, error)
}
