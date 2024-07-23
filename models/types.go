package models

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type ModelWrapper struct {
	BedrockRuntimeClient *bedrockruntime.Client
}

type LLM interface {
	Invoke() (string, error)
	Stream() (*bedrockruntime.InvokeModelWithResponseStreamOutput, error)
}

type Llama struct {
	bedrock ModelWrapper
	prompt  string
}

type Anthropic struct {
	bedrock ModelWrapper
	prompt  string
}

type StreamingOutputHandler func(ctx context.Context, part []byte) error

type ProcessingFunction func(output *bedrockruntime.InvokeModelWithResponseStreamOutput, handler StreamingOutputHandler) (any, error)

type GenericResponse interface {
	SetContent(content string)
	GetContent() string
}
