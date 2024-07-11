package models

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

func (wrapper ModelWrapper) LoadStreamingModel(modelName string, prompt string) (*bedrockruntime.InvokeModelWithResponseStreamOutput, error) {
	if modelName == anthropic {
		return wrapper.InvokeLAnthropicStream(prompt)
	} else if modelName == llama3 {
		return wrapper.InvokeLlama2Stream(prompt)
	}
	return nil, fmt.Errorf("unsupported model: %s", modelName)
}

func (wrapper ModelWrapper) LoadModel(modelName string, prompt string) (string, error) {
	if modelName == llama3 {
		return wrapper.InvokeLlama2(prompt)
	}
	return "", fmt.Errorf("unsupported model: %s", modelName)
}
