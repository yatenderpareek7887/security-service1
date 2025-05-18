package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	mysqlconfig "github.com/yatender-pareek/log-ingestor-service/src/config/my-sql-config"
	"github.com/yatender-pareek/log-ingestor-service/src/config/swagger"
	"github.com/yatender-pareek/log-ingestor-service/src/middleware"
	"github.com/yatender-pareek/log-ingestor-service/src/routes"
)

// @title Log Ingestor Service API
// @version 1.0
// @description API for ingesting and querying logs
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token
func main() {
	if err := godotenv.Load(); err != nil {
		// log.Println(".env Not found")
	}

	if err := mysqlconfig.Init(); err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	ratelimiter := middleware.NewRateLimiter(2, 5)

	r := gin.Default()
	swagger.SetupSwagger(r)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Log-Ingestor-Service!"})
	})

	basepath := r.Group(os.Getenv("BASE_PATH"))
	basepath.Use(middleware.RateLimitMiddleware(ratelimiter))
	{
		routes.SetupPublicRoutes(basepath)
	}

	basepath.Use(middleware.AuthMiddleware())
	{
		routes.SetupProtectedRoutes(basepath)
	}

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Println("Starting Log Ingestor Service on port:", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
