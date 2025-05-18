package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/yatender-pareek/log-ingestor-service/src/controllers/log-controller"
)

func SetupRouter(r *gin.RouterGroup) *gin.RouterGroup {
	r.POST("/logs", controllers.CreateLog)
	r.GET("/logs", controllers.GetAllLogs)
	r.GET("/logs/search", controllers.SearchLogs)
	r.GET("/logs/:logId", controllers.GetLogByID)
	r.DELETE("/logs/:logId", controllers.DeletelogByID)

	return r
}
