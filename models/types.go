package models

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type InvokeModelStreamingWrapper struct {
	BedrockRuntimeClient *bedrockruntime.Client
}

type StreamingOutputHandler func(ctx context.Context, part []byte) error
