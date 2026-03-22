package models

import (
	"time"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Room struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"not null"`
	Description string
	Capacity    int
	Slots       []Slot    `gorm:"foreignKey:RoomID"`
}


type Slot struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoomID    uuid.UUID `gorm:"type:uuid;index:idx_room_time"`
	StartTime time.Time `gorm:"index:idx_room_time"`
	EndTime   time.Time
	Booking   *Booking  `gorm:"foreignKey:SlotID"`
}


type Booking struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	SlotID    uuid.UUID `gorm:"type:uuid;not null"` // Уникальность ниже в SQL
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_user_bookings"`    // Индекс для "Мои брони"
	Status    string    `gorm:"default:'active'"`
	CreatedAt time.Time
}

type Schedule struct {
	ID         uuid.UUID     `gorm:"type:uuid;primaryKey"`
	RoomID     uuid.UUID     `gorm:"type:uuid;index;not null"`
	DaysOfWeek pq.Int32Array `gorm:"type:integer[]"` 
	StartTime  string        `gorm:"type:varchar(5);not null"` 
	EndTime    string        `gorm:"type:varchar(5);not null"`
}

type User struct{
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email string
	Password string
	Role string `gorm:"default:'user'"`
	Bookings []Booking `gorm:"foreignKey:UserID"`
}