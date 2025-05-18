package logingestorservice

import (
	"fmt"
	"net"
	"time"

	mysqlconfig "github.com/yatender-pareek/log-ingestor-service/src/config/my-sql-config"
	logdto "github.com/yatender-pareek/log-ingestor-service/src/dtos/log-dto"
	logDataentity "github.com/yatender-pareek/log-ingestor-service/src/models/log-data-model"
	"gorm.io/gorm"
)

type LogIngestorService struct {
}

func NewLogIngestorService() *LogIngestorService {
	return &LogIngestorService{}
}

func (s *LogIngestorService) CreateLog(dto logdto.CreateLogRequest) (*logDataentity.LogData, error) {
	if net.ParseIP(dto.IPAddress) == nil {
		return nil, fmt.Errorf("invalid IP address format")
	}

	logEntry := &logDataentity.LogData{
		Timestamp:     dto.Timestamp,
		UserID:        dto.UserID,
		IPAddress:     dto.IPAddress,
		Action:        dto.Action,
		FileName:      dto.FileName,
		DatabaseQuery: dto.DatabaseQuery,
	}

	if err := mysqlconfig.GetDB().Create(logEntry).Error; err != nil {
		return nil, fmt.Errorf("failed to save log: %v", err)
	}
	return logEntry, nil
}

func (s *LogIngestorService) GetAllLogs() ([]logDataentity.LogData, error) {
	var logs []logDataentity.LogData
	if err := mysqlconfig.GetDB().Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve logs: %v", err)
	}
	return logs, nil
}

func (s *LogIngestorService) GetLogByID(id uint64) (*logDataentity.LogData, error) {
	var logEntry logDataentity.LogData
	if err := mysqlconfig.GetDB().First(&logEntry, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("log not found")
		}
		return nil, fmt.Errorf("failed to retrieve log: %v", err)
	}
	return &logEntry, nil
}

func (s *LogIngestorService) DeleteLogByID(id uint64) error {
	result := mysqlconfig.GetDB().Where("id = ?", id).Delete(&logDataentity.LogData{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *LogIngestorService) SearchLogs(startTime, endTime *time.Time, source, userID *string) ([]logDataentity.LogData, error) {
	query := mysqlconfig.GetDB().Model(&logDataentity.LogData{})
	if startTime != nil {
		query = query.Where("timestamp >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("timestamp <= ?", *endTime)
	}
	if source != nil {
		query = query.Where("ip_address = ?", *source)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	var logs []logDataentity.LogData
	if err := query.Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to search logs: %v", err)
	}
	return logs, nil
}
