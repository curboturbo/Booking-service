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

func (mid *AuthMiddleware) auth(c *gin.Context) (string, string, bool){
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, domain.NewError(
                domain.ErrCodeUnauthorized, "Не авторизован"))
            c.Abort()
            return "","",false
        }
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
            c.JSON(http.StatusUnauthorized, domain.NewError(
                domain.ErrCodeUnauthorized, "Не авторизован"))
            c.Abort()
            return "","",false
        }

        token := parts[1]
        userID, role, err := mid.tokenProvider.VerifyToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, domain.NewError(
                domain.ErrCodeUnauthorized, "Не авторизован"))
            c.Abort()
            return "","",false
        }
        return userID, role, true
}


func (mid *AuthMiddleware) AuthenticationAdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, role, ok := mid.auth(c)
        if !ok { return }
        if (role == "admin"){
            c.Set("userID", userID)
            c.Set("role", role)
            c.Next()
        }else{
            c.JSON(http.StatusForbidden, domain.NewError(
                domain.ErrCodeForbidden, "Доступ запрещён (требуется роль admin)"))
            c.Abort()
            return
        }
    }
}


func (mid *AuthMiddleware) AuthenticationUserMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, role, ok := mid.auth(c)
        if !ok { return }
        if (role == "user"){
            c.Set("userID", userID)
            c.Set("role", role)
            c.Next()
        }else{
            c.JSON(http.StatusForbidden, domain.NewError(
                domain.ErrCodeForbidden, "Доступ запрещён (требуется роль admin)"))
            c.Abort()
            return
        }
    }
}


func (mid *AuthMiddleware) AuthenticationMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, role, ok := mid.auth(c)
        if !ok{return}
        c.Set("userID", userID)
        c.Set("role", role)
        c.Next()
    }
}