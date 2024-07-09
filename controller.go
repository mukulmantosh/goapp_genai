package main

import (
	"context"
	"goapp_genai/models"
	"log"
	"net/http"
)

func (model ModelStreamingWrapper) executeModel(w http.ResponseWriter, r *http.Request) {

	conn, err := websocketUpgrade.Upgrade(w, r, nil) // error ignored for sake of simplicity
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		return
	}

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		aiResponse, err := model.wrapper.InvokeLlama2Stream(string(msg))

		_, err = models.ProcessStreamingOutput(aiResponse, func(ctx context.Context, part []byte) error {
			err = conn.WriteMessage(msgType, part)
			if err != nil {
				log.Println("Error writing to websocket:", err)
				return nil
			}

			return nil
		})

		if err != nil {
			log.Fatal("streaming output processing error: ", err)
		}

	}
}
