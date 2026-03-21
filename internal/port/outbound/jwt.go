package port

import (
	"time"

	"github.com/google/uuid"
)

type TokenProvider interface {
	CreateToken(userID uuid.UUID, role string, duration time.Duration) (token string, err error)
	VerifyToken(token string) (userID string, role string, err error)
}