package services

import (
	"fmt"
	"log"
	"time"

	mysqlconfig "github.com/yatender-pareek/threat-analyzer-service/src/config/my-sql-config"
	threatentity "github.com/yatender-pareek/threat-analyzer-service/src/models/threat-model"
	"github.com/yatender-pareek/threat-analyzer-service/src/utility"
	"gorm.io/gorm"
)

type ThreatService struct {
}

func NewThreatService() *ThreatService {
	return &ThreatService{}
}
func (s *ThreatService) AnalyzeThreats(start *time.Time, end *time.Time) (int, error) {
	db := mysqlconfig.GetDB()
	if db == nil {
		log.Fatal("DB connection is nil!")
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error accessing raw DB connection: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}

	affectedRows, err := utility.ProcessLogs(mysqlconfig.GeSqltDB(), mysqlconfig.GetDB())

	return affectedRows, err
}

func (s *ThreatService) GetAllThreats() ([]threatentity.Threat, error) {
	var threats []threatentity.Threat
	if err := mysqlconfig.GetDB().Find(&threats).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve: %v", err)
	}
	return threats, nil
}

func (s *ThreatService) GetThreatByID(id uint64) (threatentity.Threat, error) {
	var threat threatentity.Threat
	if err := mysqlconfig.GetDB().First(&threat, id).Error; err != nil {
		return threatentity.Threat{}, err
	}
	return threat, nil
}
func (s *ThreatService) DeleteThreatByID(id uint64) error {
	result := mysqlconfig.GetDB().Where("id = ?", id).Delete(&threatentity.Threat{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *ThreatService) SearchThreats(threatType, userID, startTime, endTime string) ([]threatentity.Threat, error) {
	query := mysqlconfig.GetDB().Model(&threatentity.Threat{})
	if threatType != "" {
		query = query.Where("threat_type = ?", threatType)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if startTime != "" {
		query = query.Where("timestamp >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("timestamp <= ?", endTime)
	}

	var threats []threatentity.Threat
	if err := query.Find(&threats).Error; err != nil {
		return nil, err
	}
	return threats, nil
}
