package port

import (
	"context"
	"test-backend-1-curboturbo/internal/domain"
	"test-backend-1-curboturbo/internal/model"

	"github.com/google/uuid"
)

type StorageProvider interface {
	Create(ctx context.Context, email string, password string, role string) (models.User, error)
	GetUser(ctx context.Context, email string) (id uuid.UUID, hashpassword string,role string, err error)
	CreateAdmin(ctx context.Context, email string, password string) (error)
	
	ShowRooms(ctx context.Context) ([]domain.Room, error)
	CreateRoom(ctx context.Context, room domain.Room) (domain.Room, error)

	CreateSchedule(ctx context.Context, sched domain.Schedule) (domain.Schedule, error)
	TakeSlots(ctx context.Context, filterSlot domain.Slot) ([]domain.Slot, error)

	CreateBooking(ctx context.Context, booking domain.Booking) (domain.Booking, error)

	TakeUserBooking(ctx context.Context, userID uuid.UUID) ([]domain.Booking, error)

	CancelUserBooking(ctx context.Context, booking domain.Booking) (domain.Booking, error)

	GetAllBooking(ctx context.Context, pag domain.PaginationParams) ([]domain.Booking, error)

}