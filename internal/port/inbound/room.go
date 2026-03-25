// команды для юзера
package port

import (
	"github.com/gin-gonic/gin"
)


type RoomHandler interface{
	DisplayRooms(c *gin.Context)
	CreateRoom(c *gin.Context)
	CreateSchedule(c *gin.Context)
	TakeAvailableSlots(c *gin.Context)
	CreateReserving(c *gin.Context)
	GetInfo(c *gin.Context)
	TakeUserBookings(c *gin.Context)
	GetListOfBooking(c *gin.Context)
	CancelBooking(c *gin.Context)
}