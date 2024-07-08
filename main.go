package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/gorilla/websocket"
	"goapp_genai/aws_genai"
	"log"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "ap-south-1"
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	brc := bedrockruntime.NewFromConfig(cfg)

	invokeModel := aws_genai.InvokeModelWrapper{BedrockRuntimeClient: brc}

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {

		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			//m := strings.ToUpper(string(msg))
			message := string(msg)
			aiResponse, err := invokeModel.InvokeLlama2Stream(message)

			_, err = aws_genai.ProcessStreamingOutput(aiResponse, func(ctx context.Context, part []byte) error {
				//fmt.Println("coming..")
				//fmt.Print(string(part))
				if err = conn.WriteMessage(msgType, part); err != nil {
					return nil
				}

				return nil
			})

			if err != nil {
				log.Fatal("streaming output processing error: ", err)
			}

		}
	})

	http.ListenAndServe(":8080", nil)
}
