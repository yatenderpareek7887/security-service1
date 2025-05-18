package routes

import (
	"github.com/gin-gonic/gin"
	threatcontroller "github.com/yatender-pareek/threat-analyzer-service/src/controllers/threat-controller"
)

func SetupRouter(r *gin.RouterGroup) *gin.RouterGroup {
	r.POST("/threats/analyze", threatcontroller.AnalyzeThreats)
	r.GET("/threats", threatcontroller.GetAllThreats)
	r.GET("/threats/search", threatcontroller.SearchThreats)
	r.GET("/threats/:threatId", threatcontroller.GetThreatByID)
	r.DELETE("/threats/:threatId", threatcontroller.DeletethreatByID)

	return r
}
