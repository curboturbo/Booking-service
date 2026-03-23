package models

import (
	"time"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Room struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `gorm:"not null;unique"`
	Description string
	Capacity    int
	CreatedAt time.Time
	Slots       []Slot    `gorm:"foreignKey:RoomID"`
}


type Slot struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RoomID    uuid.UUID `gorm:"type:uuid;index:idx_room_time"`
	StartTime time.Time `gorm:"index:idx_room_time"`
	EndTime   time.Time
	Booking   *Booking  `gorm:"foreignKey:SlotID"`
}


type Booking struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SlotID    uuid.UUID `gorm:"type:uuid;not null;unique"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_user_bookings"`
	Status    string    `gorm:"default:'active'"`
	Link string 		`gorm:"default:'-'"`
	CreatedAt time.Time
}

type Schedule struct {
	ID         uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RoomID     uuid.UUID     `gorm:"type:uuid;unique;index;not null"`
	DaysOfWeek pq.Int32Array `gorm:"type:integer[]"` 
	StartTime  string        `gorm:"type:varchar(5);not null"` 
	EndTime    string        `gorm:"type:varchar(5);not null"`
}

type User struct{
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email string `gorm:"uniqueIndex"`
	Password string
	Role string `gorm:"default:'user'"`
	Bookings []Booking `gorm:"foreignKey:UserID"`
}

//CREATE UNIQUE INDEX idx_unique_active_booking ON bookings (slot_id) WHERE status = 'active';
// НЕ ЗАБЫТЬ ПОСТАВИТЬ ЧТОБЫ ПРОНИРАОВАНИЯ НЕ УЛЕТЕЛИ К ФИГАМ