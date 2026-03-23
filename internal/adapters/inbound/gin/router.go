package main_gin_router

import (
	middleware "test-backend-1-curboturbo/internal/adapters/inbound/middleware"
	auth "test-backend-1-curboturbo/internal/adapters/inbound/gin/auth"
	booking "test-backend-1-curboturbo/internal/adapters/inbound/gin/booking"
	port "test-backend-1-curboturbo/internal/port/outbound"
	service "test-backend-1-curboturbo/internal/service"
	"github.com/gin-gonic/gin"
)

func RouterInit(service service.AuthService, service_booking service.RoomService, tokenProvider port.TokenProvider)  (*gin.Engine) {
	router := gin.Default()
	authHandler := auth.NewAuthHandler(service)
	bookingHandler := booking.NewRoomHandler(service_booking)
	authMiddleware := middleware.NewAuthMiddleware(tokenProvider)

	api := router.Group("")
	{
		api.POST("/login", authHandler.Login)
		api.POST("/register",authHandler.Register)
		api.POST("/dummyLogin", authHandler.DummyLogin)
		admin_permisson := api.Group("")
		admin_permisson.Use(authMiddleware.AuthenticationAdminMiddleware())
		{
			//admin_premision -- КОГДА ВЫВОД В ПРОД
			api.POST("/rooms/create",bookingHandler.CreateRoom)
			api.POST("/rooms/:roomId/schedule/create", bookingHandler.CreateSchedule)
		}
		auth_permission := api.Group("")
		auth_permission.Use(authMiddleware.AuthenticationMiddleware())
		{
			//auth_premission -- КОГДА ВЫВОД В ПРОД
			api.GET("/rooms/list", bookingHandler.DisplayRooms)
			api.GET("/rooms/{roomId}/slots/list", bookingHandler.TakeAvailableSlots)
		}
		user_premission := api.Group("")
		user_premission.Use(authMiddleware.AuthenticationUserMiddleware())
		{
			api.POST("/bookings/create", bookingHandler.CreateReserving)
		}
	}
	return router
}