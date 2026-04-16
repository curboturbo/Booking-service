package tokenizer


import (
	"time"
	"fmt"
	"errors"
	"os"
	port "test-backend-1-curboturbo/internal/port/outbound"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)


var secretKey = []byte(os.Getenv("JWT_KEY"))

type tokenGenerator struct{}

func NewTokenGenerator() port.TokenProvider{
	return &tokenGenerator{}
}


func (t *tokenGenerator) CreateToken(userID uuid.UUID, role string, duration time.Duration) (string, error){
    claims := jwt.MapClaims{
        "userID": userID.String(),
        "role": role,
        "exp":  time.Now().Add(duration).Unix(),
        "iat":  time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    access_token, err := token.SignedString(secretKey)
    if err != nil {
        return "", err
    }
    return access_token, nil
}


func (t *tokenGenerator) VerifyToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return "", "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, okID := claims["userID"].(string)
		role, okRole := claims["role"].(string)
		if !okID || !okRole {
			return "", "", errors.New("invalid token claims")
		}
		return userID, role, nil
	}
	return "", "", errors.New("invalid token")
}