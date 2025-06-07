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

/**
* Write es un método que maneja la escritura de mensajes a través de una conexión WebSocket.
* Este método se ejecuta en un bucle infinito y tiene dos responsabilidades principales:
*
* 1. Enviar mensajes al cliente:
*   - Escucha mensajes en el canal 'send' del cliente
*   - Cuando recibe un mensaje, establece un tiempo límite de 10 segundos para escribir
*   - Si el canal está cerrado, envía un mensaje de cierre y termina
*   - Escribe el mensaje actual y cualquier mensaje en cola en la conexión WebSocket
*
* 2. Mantener la conexión viva:
*   - Utiliza un temporizador (ticker) que se activa cada 54 segundos
*   - Cuando el temporizador se activa, envía un mensaje "ping" al cliente
*   - Esto evita que la conexión se cierre por inactividad
*
* El método incluye manejo de errores y limpieza de recursos:
* - Si ocurre cualquier error al escribir, el método termina
* - Cuando el método termina, detiene el temporizador y cierra la conexión
**/
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
