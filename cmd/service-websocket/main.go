package main

import (
	"fmt"
	"log"
	"net/http"

	"fajarlaksono.github.io/laksono-api-service/app/service/websocketservice"
	"github.com/gorilla/websocket"
)

func main() {
	server := websocketservice.NewServer()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("upgrade error:", err)
			return
		}
		server.HandleWS(ws)
	})

	// http.HandleFunc("/broadcast", func(w http.ResponseWriter, r *http.Request) {
	// 	message := r.FormValue("message")
	// 	if message == "" {
	// 		http.Error(w, "Message is required", http.StatusBadRequest)
	// 		return
	// 	}
	// 	server.BroadcastMessage(message)
	// 	fmt.Fprintf(w, "Message broadcasted: %s", message)
	// })

	log.Println("WebSocket server started on :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
