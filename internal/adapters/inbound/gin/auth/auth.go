package auth

import (
	"errors"
	"net/http"
	domain "test-backend-1-curboturbo/internal/domain"
	handler "test-backend-1-curboturbo/internal/port/inbound"
	service "test-backend-1-curboturbo/internal/service"

	"github.com/gin-gonic/gin"
)

type authHandler struct{
	userAuthService service.AuthService
}


func NewAuthHandler(authService service.AuthService) handler.AuthHandler{
	return &authHandler{userAuthService: authService}
}

func (auth *authHandler) Register(c *gin.Context) {
    var req domain.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, domain.NewError(
            domain.ErrCodeInvalidRequest, 
            "invalid body: "+err.Error(),
        ))
        return
    }
    ctx := c.Request.Context()
    user, err := auth.userAuthService.Register(ctx, req.Email, req.Password, req.Role)
    if err != nil {
        if errors.Is(err, domain.ErrEmailAlreadyTaken) {
            c.JSON(http.StatusBadRequest, domain.NewError(
                domain.ErrCodeInvalidRequest,
                "Неверный запрос или email уже занят",
            ))
            return
        }
        c.JSON(http.StatusInternalServerError, domain.NewError(
            domain.ErrCodeInternalError,
            "Внутренняя ошибка сервера",
        ))
        return
    }
    c.JSON(http.StatusCreated, gin.H{
        "user": user, 
    })
}




func (auth *authHandler) Login(c *gin.Context){
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req);err!=nil{
		c.JSON(http.StatusBadRequest, domain.NewError(
        domain.ErrCodeInvalidRequest, 
        "invalid body: "+err.Error(),
        ))
		return
  	}

	ctx := c.Request.Context()
  	access_token, err := auth.userAuthService.Login(ctx, req.Email, req.Password)
  	if err != nil{
            if errors.Is(err, domain.ErrUserNotFound){
		    c.JSON(http.StatusUnauthorized, domain.NewError(
		    	domain.ErrCodeUnauthorized, "Неверные учётные данные"))
		    return
            }
            c.JSON(http.StatusInternalServerError, domain.NewError(
            domain.ErrCodeInternalError,
            "Внутренняя ошибка сервера",
        ))
        return
  	    }
	c.JSON(http.StatusOK, gin.H{
        "token": access_token,
    })
}


func (auth *authHandler) DummyLogin(c *gin.Context){
	var req domain.DummyRequest
	if err := c.ShouldBindJSON(&req);err!=nil{
		c.JSON(http.StatusBadRequest, domain.NewError(
        domain.ErrCodeInvalidRequest, 
        "invalid body: "+err.Error(),
        ))
		return
  	}
	
	ctx := c.Request.Context()
	access_token, err := auth.userAuthService.DummyLogin(ctx, req.Role)
	if err != nil{
		c.JSON(http.StatusBadRequest, domain.NewError(
			domain.ErrCodeInvalidRequest, "invalid role provided"))
		return
 	}
	c.JSON(http.StatusOK, gin.H{
        "token":access_token,
    })
}