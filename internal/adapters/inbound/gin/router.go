package main_gin_router

import (
	//middleware "test-backend-1-curboturbo/internal/adapters/inbound/middleware"
	auth "test-backend-1-curboturbo/internal/adapters/inbound/gin/auth"
	booking "test-backend-1-curboturbo/internal/adapters/inbound/gin/booking"
	port "test-backend-1-curboturbo/internal/port/outbound"
	service "test-backend-1-curboturbo/internal/service"
	"github.com/gin-gonic/gin"
)

func RouterInit(service service.AuthService, service_booking service.RoomService, tokenProvider port.TokenProvider)  (*gin.Engine) {
	router := gin.Default()
	authHandler := auth.NewAuthHandler(service)
	bookingHandler := booking.NewAuthHandler(service_booking)
	//authMiddleware := middleware.NewAuthMiddleware(tokenProvider)
	api := router.Group("")
	{
		api.POST("/login", authHandler.Login)
		api.POST("/register",authHandler.Register)
		api.POST("/dummyLogin", authHandler.DummyLogin)
		api.POST("/rooms/create",bookingHandler.CreateRoom)
		api.GET("/rooms/list", bookingHandler.DisplayRooms)
		//protected := api.Group("")
		//protected.Use(authMiddleware.AuthenticationMiddleware())

	}
	return router
}