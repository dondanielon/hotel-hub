package game

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type GameClient struct {
	connection *websocket.Conn
	send       chan []byte
	gateway    *GameGateway
}

func (client *GameClient) Read() {
	defer func() {
		client.gateway.unregister <- client
		client.connection.Close()
	}()

	client.connection.SetReadLimit(512)

	for {
		_, message, err := client.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				client.gateway.logger.Printf("WebSocket error: %v", err)
			}
			break
		}

		var wsMessage GameWebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			client.gateway.logger.Printf("Error unmarshaling message: %v", err)
			continue
		}

		client.gateway.logger.Printf("Received message: %v", wsMessage)

		client.handleMessage(wsMessage.Event, wsMessage.Payload)
		client.sendMessage(GameEvent(2), map[string]string{"message": "Hello from server"})
	}
}

func (client *GameClient) Write() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.connection.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current message
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			client.connection.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *GameClient) sendMessage(event GameEvent, payload interface{}) {
	message := GameWebSocketMessage{
		Event:   event,
		Payload: payload,
	}

	data, err := json.Marshal(message)
	if err != nil {
		client.gateway.logger.Printf("Error marshaling message: %v", err)
	}

	select {
	case client.send <- data:
	default:
		close(client.send)
	}
}

func (client *GameClient) handleMessage(event GameEvent, payload interface{}) {

}
