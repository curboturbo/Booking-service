package auth

import (
	//"errors"
	"net/http"
	domain "test-backend-1-curboturbo/internal/domain"
	handler "test-backend-1-curboturbo/internal/port/inbound"
	service "test-backend-1-curboturbo/internal/service"
	"github.com/gin-gonic/gin"
)

type roomHandler struct{
	roomService service.RoomService
}

func NewAuthHandler(roomService service.RoomService) handler.RoomHandler{
	return &roomHandler{roomService:roomService}
}

func (r *roomHandler) DisplayRooms(c *gin.Context){
	ctx := c.Request.Context()
	rooms, err := r.roomService.DisplayRooms(ctx)
	if err != nil{
		c.JSON(http.StatusInternalServerError, domain.NewError(
            domain.ErrCodeInternalError,
            "Внутренняя ошибка сервера",
        ))
	}
	c.JSON(http.StatusOK, gin.H{
        "rooms": rooms,
    })
}

func (r *roomHandler) CreateRoom(c *gin.Context){
	var req domain.RoomCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, domain.NewError(
            domain.ErrCodeInvalidRequest, 
            "неверный запрос "+err.Error(),
        ))
        return
    }

	con := c.Request.Context()
	room, err := r.roomService.CreateRoom(con, req)
	if err != nil{
		c.JSON(http.StatusInternalServerError, domain.NewError(
            domain.ErrCodeInternalError,
            "Внутренняя ошибка сервера",
        ))
	}
	c.JSON(http.StatusOK, gin.H{
		"room":room,
	})
}

