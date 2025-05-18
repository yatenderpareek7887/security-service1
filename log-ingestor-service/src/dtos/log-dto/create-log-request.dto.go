package logdto

import "time"

type CreateLogRequest struct {
	Timestamp     time.Time `json:"timestamp" validate:"required"`
	UserID        string
	IPAddress     string
	Action        string
	FileName      *string `json:"fileName"`
	DatabaseQuery *string `json:"databaseQuery"`
}
