package handler

import (
	"ais-summoner/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserByIdHandler(mongo *database.MongoDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := mongo.UserRepository().GetByID(ctx, ctx.Param("id"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}

func GetUserByEmailHandler(mongo *database.MongoDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := mongo.UserRepository().GetByEmail(ctx, ctx.Param("email"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, user)
	}
}

func GetUserListHandler(mongo *database.MongoDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := mongo.UserRepository().Find(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, users)
	}
}
