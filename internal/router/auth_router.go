package router

import (
	"ais-summoner/internal/handler"
	"ais-summoner/internal/pkg/authenticator"

	"github.com/gin-gonic/gin"
)

func NewAuthRouterV1(router *gin.Engine, auth *authenticator.Authenticator) {
	pathPrefix := "v1/auth"

	router.GET(pathPrefix+"/login", handler.LoginHandler(auth))
	router.GET(pathPrefix+"/logout", handler.LogoutHandler)
	router.GET(pathPrefix+"/callback", handler.CallbackHandler(auth))
}
