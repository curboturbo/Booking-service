package port

import (
	"context"
	"github.com/google/uuid"
)

type StorageProvider interface {
	Create(ctx context.Context, email string, password string) (error)
	GetUser(ctx context.Context, email string) (id uuid.UUID, hashpassword string,role string, err error)
	CreateAdmin(ctx context.Context, email string, password string) (error)
}