package router

import (
	"ais-summoner/internal/handler"
	"ais-summoner/internal/pkg/authenticator"

	"github.com/gin-gonic/gin"
)

func NewAuthV1Router(router *gin.Engine, auth *authenticator.Authenticator) {
	routerPrefix := "v1/auth"

	router.GET(routerPrefix+"/login", handler.LoginHandler(auth))
	router.GET(routerPrefix+"/logout", handler.LogoutHandler)
	router.GET(routerPrefix+"/callback", handler.CallbackHandler(auth))
}
