package main

import (
	"context"
	"goapp_genai/models"
	"log"
	"net/http"
)

func (m MLWrapper) executeModel(w http.ResponseWriter, r *http.Request) {
	modelName := r.URL.Query().Get("model")
	streaming := StringToBool(r.URL.Query().Get("streaming"))

	conn, err := websocketUpgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade:", err)
		return
	}
	defer conn.Close()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}
		if streaming {

			aiResponse, err := m.wrapper.LoadStreamingModel(modelName, string(msg))

			processFunc := func(ctx context.Context, part []byte) error {
				err = conn.WriteMessage(msgType, part)
				if err != nil {
					log.Println("Error writing to websocket:", err)
					return err
				}
				return nil
			}

			_, err = models.CallStreamingOutputFunction(modelName, aiResponse, processFunc)
			if err != nil {
				log.Fatal("streaming output processing error: ", err)
			}
		} else {
			modelResponse, err := m.wrapper.LoadModel(modelName, string(msg))
			err = conn.WriteMessage(msgType, []byte(modelResponse))
			if err != nil {
				log.Println("Error writing to websocket:", err)
				return
			}
		}

	}

}
