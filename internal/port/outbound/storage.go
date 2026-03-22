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
	ShowRooms(ctx context.Context) ([]models.Room, error)
	CreateRoom(ctx context.Context, room domain.Room) (models.Room, error)
}