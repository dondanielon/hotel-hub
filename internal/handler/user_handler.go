package handler

import (
	"ais-summoner/internal/database"

	"github.com/gin-gonic/gin"
)

func GetUserByIdHandler(mongo *database.MongoDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
