package middleware

import (
	"net/http"
	port "test-backend-1-curboturbo/internal/port/outbound"
	domain "test-backend-1-curboturbo/internal/domain"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct{
	tokenProvider port.TokenProvider
}

func NewAuthMiddleware(tp port.TokenProvider) *AuthMiddleware{
	return &AuthMiddleware{tokenProvider: tp}
}

func (mid *AuthMiddleware) AuthenticationMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		token,err := c.Cookie("access_token")
		if err != nil{
			c.JSON(http.StatusUnauthorized, domain.NewError(
				domain.ErrCodeUnauthorized, "unauthorized"))
			c.Abort()
			return
		}
		userID, role, err := mid.tokenProvider.VerifyToken(token)
		if err != nil{
			c.JSON(http.StatusUnauthorized, domain.NewError(
				domain.ErrCodeUnauthorized, "unauthorized"))
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}