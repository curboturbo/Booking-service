package middleware

import (
	"net/http"
	port "test-backend-1-curboturbo/internal/port/outbound"
	domain "test-backend-1-curboturbo/internal/domain"
	"github.com/gin-gonic/gin"
	"strings"
)

type AuthMiddleware struct{
	tokenProvider port.TokenProvider
}

func NewAuthMiddleware(tp port.TokenProvider) *AuthMiddleware{
	return &AuthMiddleware{tokenProvider: tp}
}

func (mid *AuthMiddleware) AuthenticationMiddleware(permission domain.Permission) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, domain.NewError(
                domain.ErrCodeUnauthorized, "missing authorization header"))
            c.Abort()
            return
        }
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
            c.JSON(http.StatusUnauthorized, domain.NewError(
                domain.ErrCodeUnauthorized, "invalid authorization format"))
            c.Abort()
            return
        }

        token := parts[1]
        userID, role, err := mid.tokenProvider.VerifyToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, domain.NewError(
                domain.ErrCodeUnauthorized, "Не авторизован"))
            c.Abort()
            return
        }
        if (permission.Role == role){
            c.Set("userID", userID)
            c.Set("role", role)
            c.Next()
        }else{
            c.JSON(http.StatusForbidden, domain.NewError(
                domain.ErrCodeForbidden, "Нет прав"))
            c.Abort()
            return
        }
    }
}