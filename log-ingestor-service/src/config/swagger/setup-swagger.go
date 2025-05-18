package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	_ "github.com/yatender-pareek/log-ingestor-service/src/docs"
)

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8082",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Log Ingestor Service API",
	Description:      "API for ingesting and querying logs",
	InfoInstanceName: "swagger",
}

func SetupSwagger(router *gin.Engine) {
	router.GET("/swagger/api/docs/*any", gs.WrapHandler(swaggerFiles.Handler))
}
