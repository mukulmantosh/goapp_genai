package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"log"
)

type Llama2Request struct {
	Prompt       string  `json:"prompt"`
	MaxGenLength int     `json:"max_gen_len,omitempty"`
	Temperature  float64 `json:"temperature,omitempty"`
}

type Llama2Response struct {
	Generation string `json:"generation"`
}

func (r Llama2Response) SetContent(content string) {
	r.Generation = content
}

func (r Llama2Response) GetContent() string {
	return r.Generation
}

func (wrapper ModelWrapper) LlamaBody(prompt string) []byte {
	body, err := json.Marshal(Llama2Request{
		Prompt:       prompt,
		MaxGenLength: 200,
		Temperature:  0.5,
	})

	if err != nil {
		log.Fatal("failed to marshal", err)
	}
	return body
}

func (wrapper ModelWrapper) InvokeLlama2Stream(prompt string) (*bedrockruntime.InvokeModelWithResponseStreamOutput, error) {

	output, err := wrapper.BedrockRuntimeClient.InvokeModelWithResponseStream(context.TODO(), &bedrockruntime.InvokeModelWithResponseStreamInput{
		ModelId:     aws.String(Llama3modelId),
		ContentType: aws.String("application/json"),
		Body:        wrapper.LlamaBody(prompt),
	})

	if err != nil {
		ProcessError(err, Llama3modelId)
	}
	return output, nil
}

func (wrapper ModelWrapper) InvokeLlama2(prompt string) (string, error) {
	modelId := "meta.llama3-70b-instruct-v1:0"

	output, err := wrapper.BedrockRuntimeClient.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelId),
		ContentType: aws.String("application/json"),
		Body:        wrapper.LlamaBody(prompt),
	})

	if err != nil {
		ProcessError(err, modelId)
	}

	var response Llama2Response
	if err := json.Unmarshal(output.Body, &response); err != nil {
		log.Fatal("failed to unmarshal", err)
	}

	return response.Generation, nil
}

func ProcessLlamaStreamingOutput(output *bedrockruntime.InvokeModelWithResponseStreamOutput, handler StreamingOutputHandler) (Llama2Response, error) {

	var combinedResult string
	resp := Llama2Response{}

	for event := range output.GetStream().Events() {
		switch v := event.(type) {
		case *types.ResponseStreamMemberChunk:

			err := json.NewDecoder(bytes.NewReader(v.Value.Bytes)).Decode(&resp)
			if err != nil {
				return resp, err
			}

			err = handler(context.Background(), []byte(resp.Generation))
			if err != nil {
				return resp, err
			}

			combinedResult += resp.GetContent()

		case *types.UnknownUnionMember:
			fmt.Println("unknown tag:", v.Tag)

		default:
			fmt.Println("union is nil or unknown type")
		}
	}

	resp.SetContent(combinedResult)
	return resp, nil

}
