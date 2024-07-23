package models

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

func (wrapper ModelWrapper) LoadStreamingModel(modelName string, prompt string) (*bedrockruntime.InvokeModelWithResponseStreamOutput, error) {
	switch modelName {
	case llama3:
		llama := Llama{LLMPrompt{wrapper, prompt}}
		response, err := llama.Stream()
		if err != nil {
			return nil, err
		}
		return response, nil
	case anthropic:
		anth := Anthropic{LLMPrompt{wrapper, prompt}}
		response, err := anth.Stream()
		if err != nil {
			return nil, err
		}
		return response, nil
	default:
		return nil, errors.New("No such model: " + modelName)
	}
}

func (wrapper ModelWrapper) LoadModel(modelName string, prompt string) (string, error) {
	switch modelName {
	case llama3:
		llama := Llama{LLMPrompt{wrapper, prompt}}
		response, err := llama.Invoke()
		if err != nil {
			return "", err
		}
		return response, nil
	case anthropic:
		anth := Anthropic{LLMPrompt{wrapper, prompt}}
		response, err := anth.Invoke()
		if err != nil {
			return "", err
		}
		return response, nil
	default:
		return "", errors.New("No such model: " + modelName)
	}
}
