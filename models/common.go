package models

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

const (
	llama3        = "llama3"
	anthropic     = "anthropic"
	Llama3modelId = "meta.llama3-70b-instruct-v1:0"
)

func CallStreamingOutputFunction(llm string, output *bedrockruntime.InvokeModelWithResponseStreamOutput, handler StreamingOutputHandler) (interface{}, error) {
	switch llm {
	case llama3:
		return ProcessLlamaStreamingOutput(output, handler)
	case anthropic:
		return ProcessAnthropicStreamingOutput(output, handler)
	default:
		return nil, fmt.Errorf("unknown llm value: %s", llm)
	}
}
