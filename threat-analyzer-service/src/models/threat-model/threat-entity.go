package threatentity

import (
	"time"
)

type Threat struct {
	ID            uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Timestamp     time.Time `json:"timestamp" gorm:"not null"`
	UserID        string    `json:"userId" gorm:"not null;type:varchar(255)"`
	IPAddress     string    `json:"ipAddress" gorm:"not null;type:varchar(255)"`
	Action        string    `json:"action" gorm:"not null;type:varchar(255)"`
	FileName      *string   `json:"fileName" gorm:"type:varchar(255)"`
	DatabaseQuery *string   `json:"databaseQuery" gorm:"type:text"`
	CreatedAt     time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"-" gorm:"autoUpdateTime"`
	ThreatType    string    `json:"threatType" gorm:"type:varchar(255)"`
	Severity      string    `json:"severity" gorm:"type:varchar(255)"`
}
