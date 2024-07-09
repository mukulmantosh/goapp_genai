package main

import (
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/gorilla/websocket"
	"goapp_genai/models"
	"log"
	"log/slog"
	"net/http"
)

var websocketUpgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ModelStreamingWrapper struct {
	wrapper models.InvokeModelStreamingWrapper
}

func main() {

	cfg := loadConfig()
	brc := bedrockruntime.NewFromConfig(cfg)
	modelWrapper := models.InvokeModelStreamingWrapper{BedrockRuntimeClient: brc}
	wrapper := ModelStreamingWrapper{modelWrapper}

	http.HandleFunc("/ws/model", wrapper.executeModel)

	slog.Info("Server Listening on", "port", "8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
