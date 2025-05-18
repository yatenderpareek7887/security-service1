package routes

import (
	"github.com/gin-gonic/gin"
	authcontroller "github.com/yatender-pareek/log-ingestor-service/src/controllers/auth-controller"
)

func SetupPublicRoutes(r *gin.RouterGroup) *gin.RouterGroup {
	r.POST("/register", authcontroller.Register)
	r.GET("/login", authcontroller.Login)
	return r
}
