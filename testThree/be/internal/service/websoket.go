package service

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	upgrader  websocket.Upgrader
	mutex     sync.Mutex
}

func NewWebSocketService() *WebSocketService {
	ws := &WebSocketService{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	// Chạy Broadcast trong một goroutine riêng
	go ws.startBroadcast()
	return ws
}

func (ws *WebSocketService) BroadcastMessage(msg []byte) {
	ws.broadcast <- msg
}

func (ws *WebSocketService) startBroadcast() {
	for msg := range ws.broadcast {
		ws.mutex.Lock()
		for client := range ws.clients {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("WebSocket write error:", err)
				client.Close()
				delete(ws.clients, client)
			}
		}
		ws.mutex.Unlock()
	}
}

func (ws *WebSocketService) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	ws.mutex.Lock()
	ws.clients[conn] = true
	ws.mutex.Unlock()

	log.Println("New WebSocket client connected")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			ws.mutex.Lock()
			delete(ws.clients, conn)
			ws.mutex.Unlock()
			break
		}
		ws.broadcast <- message
	}
}

func (ws *WebSocketService) StartBroadcast() {
	for {
		msg := <-ws.broadcast
		ws.mutex.Lock()
		for client := range ws.clients {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("WebSocket write error:", err)
				client.Close()
				delete(ws.clients, client)
			}
		}
		ws.mutex.Unlock()
	}
}
