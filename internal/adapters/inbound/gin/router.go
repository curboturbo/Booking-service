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
		api.GET("/_info", bookingHandler.GetInfo)
		api.POST("/login", authHandler.Login)
		api.POST("/register",authHandler.Register)
		api.POST("/dummyLogin", authHandler.DummyLogin)
		admin_permisson := api.Group("")
		admin_permisson.Use(authMiddleware.AuthenticationAdminMiddleware())
		{
			admin_permisson.POST("/rooms/create",bookingHandler.CreateRoom)
			admin_permisson.POST("/rooms/:roomId/schedule/create", bookingHandler.CreateSchedule)
			admin_permisson.GET("/bookings/list", bookingHandler.GetListOfBooking)
		}
		auth_permission := api.Group("")
		auth_permission.Use(authMiddleware.AuthenticationMiddleware())
		{
			auth_permission.GET("/rooms/list", bookingHandler.DisplayRooms)
			auth_permission.GET("/rooms/{roomId}/slots/list", bookingHandler.TakeAvailableSlots)
		}

		user_premission := api.Group("")
		user_premission.Use(authMiddleware.AuthenticationUserMiddleware())
		{
			user_premission.POST("/bookings/create", bookingHandler.CreateReserving)
			user_premission.GET("/bookings/my", bookingHandler.TakeUserBookings)
			user_premission.POST("/bookings/:bookingId/cancel", bookingHandler.CancelBooking)
		}
	}
	return router
}