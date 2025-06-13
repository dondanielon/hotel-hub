package main

import (
	"ais-summoner/internal/database"
	"ais-summoner/internal/game"
	"ais-summoner/internal/pkg/authenticator"
	"ais-summoner/internal/router"
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	logger := log.New(os.Stdout, "[AIS-Summoners] ", log.LstdFlags)
	logger.Println("Starting AIS Summoners server...")

	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("Error loading .env file: %v", err)
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	gameGateway := game.CreateGameGateway(database.NewMongoDB())
	go gameGateway.Run()

	rtr := gin.Default()
	rtr.Use(func(ginCtx *gin.Context) {
		ginCtx.Header("Access-Control-Allow-Origin", "*")
		ginCtx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ginCtx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if ginCtx.Request.Method == "OPTIONS" {
			ginCtx.AbortWithStatus(204)
			return
		}

		ginCtx.Next()
	})

	gob.Register(map[string]interface{}{})
	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	rtr.Use(sessions.Sessions("auth-session", store))

	rtr.GET("/health", func(ginCtx *gin.Context) {})
	rtr.GET("/ws", func(ginCtx *gin.Context) {
		cookies := ginCtx.Request.Cookies()
		logger.Printf("Cookies: %+v", cookies)
		gameGateway.HandleWebSocketConnection(ginCtx.Writer, ginCtx.Request)
	})

	router.NewAuthV1Router(rtr, auth)

	port := os.Getenv("PORT")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: rtr,
	}

	// Start server in a goroutine
	go func() {
		logger.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	handleShutdown(server, logger)
}

func handleShutdown(server *http.Server, logger *log.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Starting server shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server shutdown completed")
}
