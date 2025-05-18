package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	_ "github.com/yatender-pareek/log-ingestor-service/src/docs"
)

func SetupSwagger(router *gin.Engine) {
	router.GET("/swagger/api/docs/*any", gs.WrapHandler(swaggerFiles.Handler))
}
