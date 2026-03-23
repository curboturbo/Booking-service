package auth

import (
	//"errors"
	"errors"
	"net/http"
	"time"
	domain "test-backend-1-curboturbo/internal/domain"
	handler "test-backend-1-curboturbo/internal/port/inbound"
	service "test-backend-1-curboturbo/internal/service"
	validator "test-backend-1-curboturbo/internal/adapters/inbound/gin/validators"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)


type roomHandler struct{
	roomService service.RoomService
}


func NewRoomHandler(roomService service.RoomService) handler.RoomHandler{
	return &roomHandler{roomService:roomService}
}


func (r *roomHandler) DisplayRooms(c *gin.Context){
	ctx := c.Request.Context()
	rooms, err := r.roomService.DisplayRooms(ctx)
	if err != nil{
		c.JSON(http.StatusInternalServerError, domain.NewError(
            domain.ErrCodeInternalError,
            err.Error(),
        ))
		return
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
            err.Error(),
        ))
        return
    }

	con := c.Request.Context()
	room, err := r.roomService.CreateRoom(con, req)
	if err != nil{
		c.JSON(http.StatusInternalServerError, domain.NewError(
            domain.ErrCodeInternalError,
            err.Error(),
        ))
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"room":room,
	})
}


func (r *roomHandler) CreateSchedule(c *gin.Context) (){
	var req domain.ScheduleCreateRequest
	if err := c.ShouldBindJSON(&req);err != nil{
		c.JSON(http.StatusBadRequest, domain.NewError(
			domain.ErrCodeInvalidRequest,
			domain.ErrInvalidScheduleData.Error(),
		))
		return
	}
	if !validator.CheckDaysOfWeek(req.DaysOfWeek){
		c.JSON(http.StatusBadRequest, domain.NewError(
			domain.ErrCodeInvalidRequest,
			domain.ErrInvalidScheduleData.Error(),
		))
		return
	}
	roomIDStr := c.Param("roomId")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
    	c.JSON(http.StatusBadRequest, domain.NewError(
        	domain.ErrCodeInvalidRequest,
        	"Неверный формат roomId",
    	))
    	return
	}
	req.RoomID = roomID
	ctx := c.Request.Context()
	schedule, err := r.roomService.CreateSchedule(ctx, req)
	if err != nil{
		if errors.Is(err, domain.ErrSchedultAlreayExist){
			c.JSON(http.StatusConflict, domain.NewError(
				domain.ErrCodeScheduleExists,
				err.Error(),
			))
			return
		}
		if errors.Is(err, domain.RoomNotFound){
			c.JSON(http.StatusNotFound, domain.NewError(
				domain.ErrCodeRoomNotFound,
				err.Error(),
			))
			return
		}
		c.JSON(http.StatusInternalServerError, domain.NewError(
			domain.ErrCodeInternalError,
			err.Error(),
		))
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"schedule":schedule,
	})
}

func (r *roomHandler) TakeAvailableSlots(c *gin.Context){
	var req domain.AvailableSlotRequest
	roomIDStr := c.Param("roomId")
	roomID, err := uuid.Parse(roomIDStr)

	if err != nil{
		c.JSON(http.StatusBadRequest,domain.NewError(
			domain.ErrCodeInvalidRequest,
			domain.ErrInvalidSlotsData.Error(),
		) )
		return
	}

	date := c.Query("date")
	dateTime, err := time.Parse("2006-01-02", date)

	if err != nil{
		c.JSON(http.StatusBadRequest,domain.NewError(
			domain.ErrCodeInvalidRequest,
			domain.ErrInvalidSlotsData.Error(),
		) )
		return

	}
	
	req.Date = dateTime
	req.RoomID = roomID
	ctx := c.Request.Context()
	slots, err := r.roomService.TakeSlots(ctx, req)
	
	if err != nil{

		if errors.Is(err, domain.RoomNotFound){
			c.JSON(http.StatusNotFound, domain.NewError(
					domain.ErrCodeRoomNotFound,
					err.Error(),
			))
			return
		}
		c.JSON(http.StatusInternalServerError, domain.NewError(
			domain.ErrCodeInternalError,
			err.Error(),
		))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"slots":slots,
	})
}

func (r *roomHandler) CreateReserving(c *gin.Context){
	var req domain.CreateBookingRequest
	if err := c.ShouldBindJSON(&req);err != nil{
		c.JSON(http.StatusBadRequest, domain.NewError(
			domain.ErrCodeInvalidRequest,
			domain.ErrInvalidCreateBookingData.Error(),
		))
		return
	}
	req.UserID, _ = uuid.Parse(c.GetString("userID"))
	ctx := c.Request.Context()
	booking, err := r.roomService.ReserveSlot(ctx,req)
	if err !=nil{
		// обработка ошибок
		c.JSON(http.StatusInternalServerError, domain.NewError(
			domain.ErrCodeInternalError,
			err.Error(),
		))
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"booking":booking,
	})
}