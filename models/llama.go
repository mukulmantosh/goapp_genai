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

type StreamingOutputHandler func(ctx context.Context, part []byte) error

type Llama2Request struct {
	Prompt       string  `json:"prompt"`
	MaxGenLength int     `json:"max_gen_len,omitempty"`
	Temperature  float64 `json:"temperature,omitempty"`
}

type Llama2Response struct {
	Generation string `json:"generation"`
}

func (wrapper InvokeModelStreamingWrapper) InvokeLlama2Stream(prompt string) (*bedrockruntime.InvokeModelWithResponseStreamOutput, error) {
	modelId := "meta.llama3-70b-instruct-v1:0"

	body, err := json.Marshal(Llama2Request{
		Prompt:       prompt,
		MaxGenLength: 512,
		Temperature:  0.5,
	})

	if err != nil {
		log.Fatal("failed to marshal", err)
	}

	output, err := wrapper.BedrockRuntimeClient.InvokeModelWithResponseStream(context.TODO(), &bedrockruntime.InvokeModelWithResponseStreamInput{
		ModelId:     aws.String(modelId),
		ContentType: aws.String("application/json"),
		Body:        body,
	})

	if err != nil {
		ProcessError(err, modelId)
	}
	return output, nil
}

func ProcessStreamingOutput(output *bedrockruntime.InvokeModelWithResponseStreamOutput, handler StreamingOutputHandler) (Llama2Response, error) {

	var combinedResult string
	resp := Llama2Response{}

	for event := range output.GetStream().Events() {
		switch v := event.(type) {
		case *types.ResponseStreamMemberChunk:

			//fmt.Println("payload", string(v.Value.Bytes))

			var resp Llama2Response
			err := json.NewDecoder(bytes.NewReader(v.Value.Bytes)).Decode(&resp)
			if err != nil {
				return resp, err
			}

			err = handler(context.Background(), []byte(resp.Generation))
			if err != nil {
				return resp, err
			}

			combinedResult += resp.Generation

		case *types.UnknownUnionMember:
			fmt.Println("unknown tag:", v.Tag)

		default:
			fmt.Println("union is nil or unknown type")
		}
	}

	resp.Generation = combinedResult

	return resp, nil
}
