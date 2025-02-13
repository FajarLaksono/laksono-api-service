package websocketservice

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
	mu    sync.Mutex
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) HandleWS(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())

	s.mu.Lock()
	s.conns[ws] = true
	s.mu.Unlock()

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	defer func() {
		s.mu.Lock()
		delete(s.conns, ws)
		s.mu.Unlock()
		ws.Close()
	}()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("read error:", err)
			}
			break
		}
		fmt.Println("received message:", string(message))
		s.BroadcastMessage(string(message))
	}
}

func (s *Server) BroadcastMessage(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.conns {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println("write error:", err)
			conn.Close()
			delete(s.conns, conn)
		}
	}
}
