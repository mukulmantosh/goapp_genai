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

type Claude3Request struct {
	AnthropicVersion string    `json:"anthropic_version"`
	MaxTokens        int       `json:"max_tokens"`
	Messages         []Message `json:"messages"`
	Temperature      float64   `json:"temperature,omitempty"`
	TopP             float64   `json:"top_p,omitempty"`
	TopK             int       `json:"top_k,omitempty"`
	StopSequences    []string  `json:"stop_sequences,omitempty"`
	SystemPrompt     string    `json:"system,omitempty"`
}
type Content struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}
type Message struct {
	Role    string    `json:"role,omitempty"`
	Content []Content `json:"content,omitempty"`
}

type Claude3Response struct {
	ID              string            `json:"id,omitempty"`
	Model           string            `json:"model,omitempty"`
	Type            string            `json:"type,omitempty"`
	Role            string            `json:"role,omitempty"`
	ResponseContent []ResponseContent `json:"content,omitempty"`
	StopReason      string            `json:"stop_reason,omitempty"`
	StopSequence    string            `json:"stop_sequence,omitempty"`
	Usage           Usage             `json:"usage,omitempty"`
}
type ResponseContent struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}
type Usage struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
}

type PartialResponse struct {
	Type    string                 `json:"type"`
	Message PartialResponseMessage `json:"message,omitempty"`
	Index   int                    `json:"index,omitempty"`
	Delta   Delta                  `json:"delta,omitempty"`
	Usage   PartialResponseUsage   `json:"usage,omitempty"`
}

type PartialResponseMessage struct {
	ID           string               `json:"id,omitempty"`
	Type         string               `json:"type,omitempty"`
	Role         string               `json:"role,omitempty"`
	Content      []interface{}        `json:"content,omitempty"`
	Model        string               `json:"model,omitempty"`
	StopReason   string               `json:"stop_reason,omitempty"`
	StopSequence interface{}          `json:"stop_sequence,omitempty"`
	Usage        PartialResponseUsage `json:"usage,omitempty"`
}

type PartialResponseUsage struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
}

type Delta struct {
	Type       string `json:"type,omitempty"`
	Text       string `json:"text,omitempty"`
	StopReason string `json:"stop_reason,omitempty"`
}

const partialResponseTypeContentBlockDelta = "content_block_delta"
const partialResponseTypeMessageStart = "message_start"
const partialResponseTypeMessageDelta = "message_delta"

const contentTypeText = "text"

func (r Claude3Response) SetContent(content string) {
	r.ResponseContent[0].Text = content
}
func (r Claude3Response) GetContent() string {
	return r.ResponseContent[0].Text
}

func (wrapper Anthropic) AnthropicBody(prompt string) []byte {
	payload := Claude3Request{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        200,
		Messages: []Message{
			{
				Role: "user",
				Content: []Content{
					{
						Type: "text",
						Text: prompt,
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func (wrapper Anthropic) Stream() (*bedrockruntime.InvokeModelWithResponseStreamOutput, error) {
	body := wrapper.AnthropicBody(wrapper.prompt)

	output, err := wrapper.bedrock.BedrockRuntimeClient.InvokeModelWithResponseStream(context.TODO(), &bedrockruntime.InvokeModelWithResponseStreamInput{
		ModelId:     aws.String(claudeV3ModelID),
		ContentType: aws.String("application/json"),
		Body:        body,
	})

	if err != nil {
		ProcessError(err, claudeV3ModelID)
	}
	return output, nil

}

func (wrapper Anthropic) Invoke() (string, error) {
	body := wrapper.AnthropicBody(wrapper.prompt)

	output, err := wrapper.bedrock.BedrockRuntimeClient.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(claudeV3ModelID),
		ContentType: aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		log.Fatal(err)
	}

	var resp Claude3Response

	err = json.Unmarshal(output.Body, &resp)

	if err != nil {
		log.Fatal(err)
	}

	return resp.ResponseContent[0].Text, nil
}

func ProcessAnthropicStreamingOutput(output *bedrockruntime.InvokeModelWithResponseStreamOutput, handler StreamingOutputHandler) error {

	var combinedResult string
	resp := Claude3Response{
		Type:            "message",
		Role:            "assistant",
		Model:           claudeV3ModelID,
		ResponseContent: []ResponseContent{{Type: contentTypeText}}}

	for event := range output.GetStream().Events() {
		switch v := event.(type) {
		case *types.ResponseStreamMemberChunk:

			var pr PartialResponse
			err := json.NewDecoder(bytes.NewReader(v.Value.Bytes)).Decode(&pr)
			if err != nil {
				return err
			}

			if pr.Type == partialResponseTypeContentBlockDelta {
				handler(context.Background(), []byte(pr.Delta.Text))
				combinedResult += pr.Delta.Text
			} else if pr.Type == partialResponseTypeMessageStart {
				resp.ID = pr.Message.ID
				resp.Usage.InputTokens = pr.Message.Usage.InputTokens
			} else if pr.Type == partialResponseTypeMessageDelta {
				resp.StopReason = pr.Delta.StopReason
				resp.Usage.OutputTokens = pr.Message.Usage.OutputTokens
			}

		case *types.UnknownUnionMember:
			return fmt.Errorf("unknown tag: %s", v.Tag)

		default:
			return fmt.Errorf("union is nil or unknown type")
		}
	}
	return nil
}
