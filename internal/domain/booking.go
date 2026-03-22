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