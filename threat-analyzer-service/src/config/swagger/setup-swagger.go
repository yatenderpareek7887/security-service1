package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	_ "github.com/yatender-pareek/threat-analyzer-service/src/docs"
)

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8083",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Threat Analyzer Service API",
	Description:      "API for threat analyzer",
	InfoInstanceName: "swagger",
}

func SetupSwagger(router *gin.Engine) {
	router.GET("/swagger/api/docs/*any", gs.WrapHandler(swaggerFiles.Handler))
}
