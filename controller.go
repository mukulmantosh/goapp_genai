package main

import (
	"context"
	"goapp_genai/models"
	"log"
	"net/http"
)

func (model ModelStreamingWrapper) executeModel(w http.ResponseWriter, r *http.Request) {

	conn, _ := websocketUpgrade.Upgrade(w, r, nil) // error ignored for sake of simplicity

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		message := string(msg)

		aiResponse, err := model.wrapper.InvokeLlama2Stream(message)

		_, err = models.ProcessStreamingOutput(aiResponse, func(ctx context.Context, part []byte) error {

			if err = conn.WriteMessage(msgType, part); err != nil {
				log.Println(err)
				return nil
			}

			return nil
		})

		if err != nil {
			log.Fatal("streaming output processing error: ", err)
		}

	}
}
