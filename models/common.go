package models

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

const (
	llama3          = "llama3"
	anthropic       = "anthropic"
	Llama3modelId   = "meta.llama3-70b-instruct-v1:0"
	claudeV3ModelID = "anthropic.claude-3-haiku-20240307-v1:0"
)

func CallStreamingOutputFunction(llm string, output *bedrockruntime.InvokeModelWithResponseStreamOutput, handler StreamingOutputHandler) error {
	switch llm {
	case llama3:
		err := ProcessLlamaStreamingOutput(output, handler)
		if err != nil {
			return err
		}
	case anthropic:
		err := ProcessAnthropicStreamingOutput(output, handler)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown llm value: %s", llm)
	}
	return nil
}
