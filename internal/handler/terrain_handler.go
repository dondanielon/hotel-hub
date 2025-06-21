package handler

import (
	"ais-summoner/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTerrainByIdHandler(mongodb *database.MongoDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		terrain, err := mongodb.TerrainRepository().GetByID(ctx, ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, terrain)
	}
}
