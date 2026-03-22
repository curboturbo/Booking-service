package service

import (
	"context"
	"errors"
	"fmt"
	domain "test-backend-1-curboturbo/internal/domain"
	model "test-backend-1-curboturbo/internal/model"
	port "test-backend-1-curboturbo/internal/port/outbound"
	"time"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var tokenLiving = 2*time.Hour

type AuthService interface{
	Register(ctx context.Context, email string, password string, role string) (model.User, error)
	Login(ctx context.Context, email string, password string) (accessToken string, err error)
	DummyLogin(ctx context.Context, role string) (accessToken string, err error)
}


type authService struct{
	storage port.StorageProvider
	tokenGen port.TokenProvider
}

func NewAuthService(storage port.StorageProvider, tokenGen port.TokenProvider) AuthService{
	return &authService{storage: storage, tokenGen: tokenGen}
}

func (s *authService) Register(ctx context.Context, email string, password string, role string) (model.User, error){
	bytes, er := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost);if er != nil{return model.User{}, er}
	user, err := s.storage.Create(ctx, email, string(bytes), role)
	if err != nil{return model.User{}, err}
	return user, nil
}


func (s *authService) Login(ctx context.Context, email string, password string) (string, error){
	userID, hash, role, err := s.storage.GetUser(ctx, email)
	if err != nil {return "",err}
	if er := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password)); er !=nil{
		return "", er
	}
	token, err := s.tokenGen.CreateToken(userID, role, tokenLiving)
	if err != nil{return "",err}
	return token, nil
}

func (s *authService) DummyLogin(ctx context.Context, role string) (string, error) {
    dummy := domain.NewGenDummy()
    var targetID uuid.UUID
    switch role {
    case "admin":
        targetID = dummy.AdminID
    case "user":
        targetID = dummy.UserID
    default:
        return "", errors.New("no available role")
    }
    token, err := s.tokenGen.CreateToken(targetID, role, tokenLiving)
    if err != nil {
        return "", fmt.Errorf("failed to generate dummy token: %w", err)
    }
    return token, nil
}