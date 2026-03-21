package domain

import "github.com/google/uuid"


type RegisterRequest struct {
	Email string `json:"email" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type DummyRequest struct {
	Role string `json:"role" binding:"required"`
}


type Dummy struct {
    UserID uuid.UUID
	AdminID uuid.UUID
}

func NewGenDummy() (Dummy){
	return Dummy{
		UserID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		AdminID:uuid.MustParse("a1b2c3d4-e5f6-47a8-b9c0-d1e2f3a4b5c6"),
	}
}

