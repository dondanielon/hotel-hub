package handler

import (
	"ais-summoner/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTerrainByIdHandler(mongodb *database.MongoDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"id": ctx.Param("id"),
		})
	}
}
