package domain
import "github.com/google/uuid"
import "time"

type RoomCreateRequest struct{
	Name string `json:"name" binding:"required"`
	Description string `json:"description"`
	Capacity int `json:"capacity"`
}

type Room struct {
	ID          uuid.UUID
	Name        string
	Description string
	Capacity    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}


type ScheduleCreateRequest struct{
	DaysOfWeek []int `json:"daysOfWeek"`
	StartTime string `json:"startTime"`
	EndTime string `json:"endTime"`
	RoomID uuid.UUID
}


type Schedule struct{
	ID uuid.UUID
	RoomID uuid.UUID
	DaysOfWeek []int
	StartTime string
	EndTime string
}

type AvailableSlotRequest struct{
	RoomID uuid.UUID `json:"roomID"`
	Date time.Time `json:"date"`
}

type Slot struct{
	RoomId uuid.UUID
	StartTime time.Time
	EndTime time.Time
}


type CreateBookingRequest struct{
	SlotID uuid.UUID `json:"slotId"`
	CreateConferenceLink bool `json:"createConferenceLink"`
	UserID uuid.UUID
}


type Booking struct{
	ID uuid.UUID
	SlotID uuid.UUID
	UserID uuid.UUID
	Status string
	Link string
}


type RequestCancelBooking struct{
	UserID uuid.UUID
	BookingID uuid.UUID
}


type PaginationParams struct {
    Page     int `form:"page" binding:"omitempty,numeric,min=1"`
    PageSize int `form:"pageSize" binding:"omitempty,numeric,min=1,max=100"`
}