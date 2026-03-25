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
		if errors.Is(err, domain.ErrSlotAlreadyTaken){
			c.JSON(http.StatusConflict,domain.NewError(
				domain.ErrCodeSlotAlreadyBooked,
				err.Error(),
			))
		}
		if errors.Is(err, domain.ErrSlotDosntExist){
			c.JSON(http.StatusNotFound, domain.NewError(
				domain.ErrCodeSlotNotFound,
				err.Error(),
			))
		}
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

func (r *roomHandler) GetListOfBooking(c *gin.Context){
	var params domain.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest,domain.NewError(
			domain.ErrCodeInvalidRequest,
			domain.ErrInvalidSlotsData.Error(),
		) )
		return
	}
	if params.Page == 0 {
        params.Page = 1
    }
    if params.PageSize == 0 {
        params.PageSize = 20
    }
	ctx := c.Request.Context()
	bookings, err := r.roomService.GetAllBooking(ctx, params)
	if err != nil{
		c.JSON(http.StatusInternalServerError, domain.NewError(
			domain.ErrCodeInternalError,
			err.Error(),
		))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"bookings":bookings,
	})

}


func (r *roomHandler) TakeUserBookings(c *gin.Context){
	UserID, _ := uuid.Parse(c.GetString("userID"))
	ctx := c.Request.Context()
	bookings, err := r.roomService.TakeUserBooking(ctx, UserID)
	if err != nil{
		c.JSON(http.StatusInternalServerError, domain.NewError(
			domain.ErrCodeInternalError,
			err.Error(),
		))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"future_bookings": bookings,
	})
}



func (r *roomHandler) CancelBooking(c *gin.Context){
	UserID, err1 := uuid.Parse(c.GetString("userID"))
	bookingIDStr := c.Param("bookingId")
	bookingID, err2 := uuid.Parse(bookingIDStr)
	if err1!=nil || err2!=nil{
		c.JSON(http.StatusBadRequest,domain.NewError(
			domain.ErrCodeInvalidRequest,
			domain.ErrInvalidSlotsData.Error(),
		) )
		return
	}
	ctx := c.Request.Context()
	booking, err := r.roomService.CancelUserBooking(ctx,domain.RequestCancelBooking{UserID:UserID, BookingID:bookingID})
	if err != nil{
		if errors.Is(err, domain.ErrBookingNotFound){
			c.JSON(http.StatusNotFound, domain.NewError(
				domain.ErrCodeBookingNotFound,
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
		"booking":booking,
	})
}



func (r *roomHandler) GetInfo(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"message":"API is available",
	})
}
