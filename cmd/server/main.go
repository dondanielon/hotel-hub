package main

import (
	"ais-summoner/internal/game"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := log.New(os.Stdout, "[AIS-Summoners] ", log.LstdFlags)
	logger.Println("Starting AIS Summoners server...")

	gameGateway := game.CreateGameGateway()
	go gameGateway.Run()

	router := gin.Default()
	router.Use(func(ginCtx *gin.Context) {
		ginCtx.Header("Access-Control-Allow-Origin", "*")
		ginCtx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ginCtx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if ginCtx.Request.Method == "OPTIONS" {
			ginCtx.AbortWithStatus(204)
			return
		}

		ginCtx.Next()
	})

	router.GET("/health", func(ginCtx *gin.Context) {})
	router.GET("/ws", func(ginCtx *gin.Context) {
		cookies := ginCtx.Request.Cookies()
		headers := ginCtx.Request.Header.Values("x-custom")
		logger.Printf("Cookies: %+v", cookies)
		logger.Printf("Headers: %+v", headers[0])
		gameGateway.HandleWebSocketConnection(ginCtx.Writer, ginCtx.Request)
	})

	port := "8080"
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
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
