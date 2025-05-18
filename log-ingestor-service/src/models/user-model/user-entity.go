package userentity

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the database
type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string         `gorm:"type:varchar(255);unique;not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"password"`
	Email     string         `gorm:"type:varchar(255);not null" json:"email"`
	CreatedAt time.Time      `json:"-" gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	UpdatedAt time.Time      `json:"-" gorm:"autoUpdateTime,omitempty"`
}
