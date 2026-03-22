package main_gin_router

import (
	middleware "test-backend-1-curboturbo/internal/adapters/inbound/middleware"
	auth "test-backend-1-curboturbo/internal/adapters/inbound/gin/http"
	port "test-backend-1-curboturbo/internal/port/outbound"
	service "test-backend-1-curboturbo/internal/service"
	"github.com/gin-gonic/gin"
)

func RouterInit(service service.AuthService, tokenProvider port.TokenProvider)  (*gin.Engine) {
	router := gin.Default()
	authHandler := auth.NewAuthHandler(service)
	authMiddleware := middleware.NewAuthMiddleware(tokenProvider)
	api := router.Group("")
	{
		api.POST("/login", authHandler.Login)
		api.POST("/register",authHandler.Register)
		api.POST("/dummyLogin", authHandler.DummyLogin)
		protected := api.Group("")
		protected.Use(authMiddleware.AuthenticationMiddleware())
		{
			protected.POST("/logout", authHandler.Logout)
		}
	}
	return router
}