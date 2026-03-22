package auth


import (
	"net/http"
	handler "test-backend-1-curboturbo/internal/port/inbound"
	service "test-backend-1-curboturbo/internal/service"
	domain "test-backend-1-curboturbo/internal/domain"
	"github.com/gin-gonic/gin"
)

type authHandler struct{
	userAuthService service.AuthService
}


func NewAuthHandler(authService service.AuthService) handler.AuthHandler{
	return &authHandler{userAuthService: authService}
}

func (auth *authHandler) Register(c *gin.Context){
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req);err!=nil{
		c.JSON(http.StatusBadRequest, domain.NewError(
            domain.ErrCodeInvalidRequest, 
            "invalid body: "+err.Error(),
        ))
		return
  }

  ctx := c.Request.Context()
  if err := auth.userAuthService.Register(ctx, req.Email, req.Password);err !=nil{
	c.JSON(http.StatusInternalServerError, domain.NewError(
            domain.ErrCodeInternalError, 
            "registration failed",
        ))
	return
  }

  c.JSON(http.StatusCreated, gin.H{
	"message": "user registered successfully"})
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
		c.JSON(http.StatusUnauthorized, domain.NewError(
			domain.ErrCodeUnauthorized, "invalid email or password"))
		return
  	}
	c.SetCookie(
        "access_token", 
        access_token, 
        3600*6,
        "/", 
        "", 
        false,
        true,
    )
	c.JSON(http.StatusOK, gin.H{
        "message": "successfully logged in",
    })
}

func (auth *authHandler) Logout(c *gin.Context){
	val, ok:= c.Get("userID")
	roleVal, okVal := c.Get("role")
	if !(ok && okVal){
		c.JSON(http.StatusUnauthorized, domain.NewError(
			domain.ErrCodeUnauthorized, "unauthorized"))
        return
	}
	userID, ok1 := val.(string)
	role, ok2 := roleVal.(string)
	if !(ok1 && ok2){
		c.JSON(http.StatusInternalServerError, domain.NewError(
			domain.ErrCodeInternalError, "internal context error"))
        return
	}
	err := auth.userAuthService.Logout(c.Request.Context(), userID, role)
	if err != nil{
		c.JSON(http.StatusInternalServerError, domain.NewError(
			domain.ErrCodeInternalError, "failed to logout"))
        return
	}
	c.SetCookie("access_token", "", -1, "/", "", false, true)
    c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
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
	c.SetCookie(
        "access_token",
        access_token, 
        3600*6,
        "/",
        "",
        false,
        true,
    )
	c.JSON(http.StatusOK, gin.H{
        "message": "successfully logged in",
    })
}