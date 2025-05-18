// Package controllers provides HTTP handlers for log-related operations
package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	logdto "github.com/yatender-pareek/log-ingestor-service/src/dtos/log-dto"
	genricerror "github.com/yatender-pareek/log-ingestor-service/src/genric_error"
	logingestorservice "github.com/yatender-pareek/log-ingestor-service/src/services/log-ingestor-service"
	"gorm.io/gorm"
)

var (
	validate   = validator.New()
	LogService *logingestorservice.LogIngestorService
)

func InitControllers() {
	// if service == nil {
	// 	panic("LogService cannot be nil")
	// }
	if LogService != nil {
		panic("LogService already initialized")
	}
	if validate != nil {
		panic("Validator already initialized")
	}

	validate = validator.New()
	validate.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		for _, r := range value {
			if !unicode.IsSpace(r) {
				return true
			}
		}
		return false
	})
	log.Printf("------------ar 44----------")
	LogService = logingestorservice.NewLogIngestorService()
}

// CreateLog godoc
// @Summary Create a new log
// @Description Create a new log record with provided details
// @Tags Logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param log body logdto.CreateLogRequest true "Log to create"
// @Success 201 {object} logdto.CreateLogRequest
// @Failure 400 {object} genricerror.ErrorResponse
// @Failure 500 {object} genricerror.ErrorResponse
// @Router /api/logs [post]
func CreateLog(c *gin.Context) {
	var logDto logdto.CreateLogRequest
	if err := c.ShouldBindJSON(&logDto); err != nil {
		c.JSON(http.StatusBadRequest, genricerror.ErrorResponse{Message: err.Error()})
		return
	}

	if err := validate.Struct(logDto); err != nil {
		c.JSON(http.StatusBadRequest, genricerror.ErrorResponse{Message: "Validation failed: " + err.Error()})
		return
	}

	createdLog, err := LogService.CreateLog(logDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, genricerror.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdLog)
}

// GetAllLogs godoc
// @Summary Get all logs
// @Description Get list of all logs
// @Tags Logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} logdto.CreateLogRequest
// @Failure 500 {object} genricerror.ErrorResponse
// @Router /api/logs [get]
func GetAllLogs(c *gin.Context) {

	logs, err := LogService.GetAllLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, genricerror.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// GetLogByID godoc
// @Summary Get a specific log
// @Description Retrieve a log by its ID
// @Tags Logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param logId path int true "Log ID"
// @Success 200 {object} logdto.CreateLogRequest
// @Failure 400 {object} genricerror.ErrorResponse
// @Failure 404 {object} genricerror.ErrorResponse
// @Failure 500 {object} genricerror.ErrorResponse
// @Router /api/logs/{logId} [get]
func GetLogByID(c *gin.Context) {
	logIDStr := c.Param("logId")
	logID, err := strconv.ParseUint(logIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, genricerror.ErrorResponse{Message: "Invalid log ID"})
		return
	}

	logEntry, err := LogService.GetLogByID(logID)
	if err != nil {
		if err.Error() == "log not found" {
			c.JSON(http.StatusNotFound, genricerror.ErrorResponse{Message: "Log not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, genricerror.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, logEntry)
}

// DeleteLogByID godoc
// @Summary Delete a specific log
// @Description Deletes a log by its ID
// @Tags Logs
// @Produce json
// @Security BearerAuth
// @Param logId path int true "Log ID"
// @Success 204 {object} nil "log deleted successfully"
// @Failure 404 {object} map[string]string "log not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/logs/{logId} [delete]
func DeletelogByID(c *gin.Context) {
	logIDStr := c.Param("logId")
	logID, err := strconv.ParseUint(logIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, genricerror.ErrorResponse{Message: "Invalid log ID"})
		return
	}
	err = LogService.DeleteLogByID(logID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// SearchLogs godoc
// @Summary Search logs
// @Description Retrieve logs based on time range, source, or user
// @Tags Logs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_time query string false "Start time (RFC3339)"
// @Param end_time query string false "End time (RFC3339)"
// @Param source query string false "Source IP address"
// @Param user_id query string false "User ID"
// @Success 200 {array} logdto.CreateLogRequest
// @Failure 400 {object} genricerror.ErrorResponse
// @Failure 500 {object} genricerror.ErrorResponse
// @Router /api/logs/search [get]
func SearchLogs(c *gin.Context) {
	var startTime, endTime *time.Time
	var source, userID *string

	if start := c.Query("start_time"); start != "" {
		t, err := time.Parse(time.RFC3339, start)
		if err != nil {
			c.JSON(http.StatusBadRequest, genricerror.ErrorResponse{Message: "Invalid start_time format"})
			return
		}
		startTime = &t
	}

	if end := c.Query("end_time"); end != "" {
		t, err := time.Parse(time.RFC3339, end)
		if err != nil {
			c.JSON(http.StatusBadRequest, genricerror.ErrorResponse{Message: "Invalid end_time format"})
			return
		}
		endTime = &t
	}

	if src := c.Query("source"); src != "" {
		source = &src
	}

	if uid := c.Query("user_id"); uid != "" {
		userID = &uid
	}

	logs, err := LogService.SearchLogs(startTime, endTime, source, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, genricerror.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}
