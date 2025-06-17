package router

import (
	"ais-summoner/internal/database"
	"ais-summoner/internal/handler"

	"github.com/gin-gonic/gin"
)

func NewUserRouterV1(router *gin.Engine, mongodb *database.MongoDB) {
	pathPrefix := "/v1/user"

	router.GET(pathPrefix+"/:id", handler.GetUserByIdHandler(mongodb))
}
