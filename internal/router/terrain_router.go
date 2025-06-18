package router

import (
	"ais-summoner/internal/database"
	"ais-summoner/internal/handler"

	"github.com/gin-gonic/gin"
)

func NewTerrainRouterV1(router *gin.Engine, mongodb *database.MongoDB) {
	pathPrefix := "/v1/terrain"

	router.GET(pathPrefix+"/:id", handler.GetTerrainByIdHandler(mongodb))
}
