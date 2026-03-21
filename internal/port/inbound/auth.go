package port

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	Logout(c *gin.Context)
	DummyLogin(c *gin.Context)
}
