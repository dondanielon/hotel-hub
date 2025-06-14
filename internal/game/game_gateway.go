package game

import (
	"ais-summoner/internal/database"

	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type GameGateway struct {
	clients    map[*GameClient]bool
	mongodb    *database.MongoDB
	logger     *log.Logger
	mutex      sync.RWMutex
	redis      *database.Redis
	register   chan *GameClient
	unregister chan *GameClient
	upgrader   websocket.Upgrader
}

func NewGameGateway(mongodb *database.MongoDB, cache *database.Redis) *GameGateway {
	return &GameGateway{
		clients: make(map[*GameClient]bool),
		mongodb: mongodb,
		logger:  log.New(log.Writer(), "[GameGateway] ", log.LstdFlags),
		redis:   cache,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (gateway *GameGateway) Run() {
	for {
		select {
		case client := <-gateway.register:
			gateway.mutex.Lock()
			gateway.clients[client] = true
			gateway.mutex.Unlock()
			gateway.logger.Printf("Client registered")

		case client := <-gateway.unregister:
			gateway.mutex.Lock()
			_, exists := gateway.clients[client]
			if exists {
				delete(gateway.clients, client)
				close(client.send)
			}
			gateway.mutex.Unlock()
			gateway.logger.Printf("Client unregistered")
		}
	}
}

func (gateway *GameGateway) HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := gateway.upgrader.Upgrade(w, r, nil)
	if err != nil {
		gateway.logger.Printf("Error upgrading connection: %v", err)
	}

	client := &GameClient{
		connection: conn,
		send:       make(chan []byte, 256),
		gateway:    gateway,
	}

	gateway.clients[client] = true

	go client.Read()
	go client.Write()
}
