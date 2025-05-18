package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	threatanalyzerresquest "github.com/yatender-pareek/threat-analyzer-service/src/dto/threat-analyzer-resquest"
	services "github.com/yatender-pareek/threat-analyzer-service/src/services/threat-service"
	"gorm.io/gorm"
)

// AnalyzeThreatRequest represents the request body for the analyze endpoint
type AnalyzeThreatRequest struct {
	Start time.Time `json:"start" binding:"required" example:"2025-03-26T00:00:00Z" format:"date-time"`
	End   time.Time `json:"end" binding:"required" example:"2025-03-27T00:00:00Z" format:"date-time"`
}

var threatService *services.ThreatService

func InitController() {
	threatService = services.NewThreatService()
}

// AnalyzeThreats godoc
// @Summary Analyze logs for threats
// @Description Analyzes logs within the specified time range and detects threats
// @Tags Threats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AnalyzeThreatRequest true "Start and end time for log analysis"
// @Success 200 {array} threatentity.Threat "List of detected threats"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/threats/analyze [post]
func AnalyzeThreats(c *gin.Context) {
	var req threatanalyzerresquest.AnalyzeThreatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	defaultStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	defaultEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	startTime := defaultStart
	if req.StartTime != nil {
		startTime = *req.StartTime
	}

	endTime := defaultEnd
	if req.EndTime != nil {
		endTime = *req.EndTime
	}

	if !endTime.After(startTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end time must be after start time"})
		return
	}

	threatsResult, err := threatService.AnalyzeThreats(req.StartTime, req.EndTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if threatsResult == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No threats detected"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Threat %d detected!!", threatsResult)})
}

// GetAllThreats godoc
// @Summary Retrieve all threats
// @Description Fetches all detected threats from the database
// @Tags Threats
// @Produce json
// @Security BearerAuth
// @Success 200 {array} threatentity.Threat "List of threats"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/threats [get]
func GetAllThreats(c *gin.Context) {
	threats, err := threatService.GetAllThreats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, threats)
}

// GetThreatByID godoc
// @Summary Retrieve a specific threat
// @Description Fetches a threat by its ID
// @Tags Threats
// @Produce json
// @Security BearerAuth
// @Param threatId path int true "Threat ID"
// @Success 200 {object} threatentity.Threat "Threat details"
// @Failure 404 {object} map[string]string "Threat not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/threats/{threatId} [get]
func GetThreatByID(c *gin.Context) {
	threatIDStr := c.Param("threatId")
	threatID, err := strconv.ParseUint(threatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in parse threat id"})
		return
	}
	threat, err := threatService.GetThreatByID(threatID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Threat not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, threat)
}

// DeletethreatByID godoc
// @Summary Delete a specific threat
// @Description Deletes a threat by its ID
// @Tags Threats
// @Produce json
// @Security BearerAuth
// @Param threatId path int true "threat ID"
// @Success 204 {object} nil "threat deleted successfully"
// @Failure 404 {object} map[string]string "threat not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/threats/{threatId} [delete]
func DeletethreatByID(c *gin.Context) {
	threatIDStr := c.Param("threatId")
	threatID, err := strconv.ParseUint(threatIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in parse threat id"})
		return
	}
	err = threatService.DeleteThreatByID(threatID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "threat not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

}

// SearchThreats godoc
// @Summary Search threats
// @Description Searches threats by type, user, or time range
// @Tags Threats
// @Produce json
// @Security BearerAuth
// @Param type query string false "Threat type"
// @Param user query string false "User ID"
// @Param startTime query string false "Start time (RFC3339)" format:"date-time"
// @Param endTime query string false "End time (RFC3339)" format:"date-time"
// @Success 200 {array} threatentity.Threat "List of matching threats"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/threats/search [get]
func SearchThreats(c *gin.Context) {
	threatType := c.Query("type")
	userID := c.Query("user")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	threats, err := threatService.SearchThreats(threatType, userID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, threats)
}
