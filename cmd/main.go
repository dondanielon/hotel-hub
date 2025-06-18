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
	loadEnvVariables(logger)

	mongodb := database.NewMongoDB()
	gateway := game.NewGameGateway(mongodb, database.NewRedis())
	go gateway.Run()

	gob.Register(map[string]interface{}{})
	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))

	ginRouter := gin.Default()
	ginRouter.Use(func(ginCtx *gin.Context) {
		ginCtx.Header("Access-Control-Allow-Origin", "*")
		ginCtx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ginCtx.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if ginCtx.Request.Method == "OPTIONS" {
			ginCtx.AbortWithStatus(204)
			return
		}

		ginCtx.Next()
	})
	ginRouter.Use(sessions.Sessions("auth-session", store))
	ginRouter.GET("/health", func(ginCtx *gin.Context) {
		session := sessions.Default(ginCtx)
		log.Printf("profile: %v", session.Get("profile"))
		log.Printf("access_token: %v", session.Get("access_token"))
	})
	ginRouter.GET("/version", func(ginCtx *gin.Context) {})
	ginRouter.GET("/ws", func(ginCtx *gin.Context) {
		gateway.HandleWebSocketConnection(ginCtx.Writer, ginCtx.Request)
	})

	auth, err := authenticator.NewAuthenticator()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	router.NewAuthRouterV1(ginRouter, auth)
	router.NewUserRouterV1(ginRouter, mongodb)
	router.NewTerrainRouterV1(ginRouter, mongodb)

	port := os.Getenv("PORT")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: ginRouter,
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

func loadEnvVariables(logger *log.Logger) {
	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("Error loading .env file: %v", err)
	}
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
