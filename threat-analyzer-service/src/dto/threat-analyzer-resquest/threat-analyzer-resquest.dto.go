package threatanalyzerresquest

import "time"

type AnalyzeThreatRequest struct {
	StartTime *time.Time `json:"startTime" validate:"omitempty"`
	EndTime   *time.Time `json:"endTime" validate:"omitempty"`
}

type SearchThreatRequest struct {
	Type      *string    `json:"type" validate:"omitempty,notblank"`
	UserID    *string    `json:"userId" validate:"omitempty,notblank"`
	StartTime *time.Time `json:"startTime" validate:"omitempty"`
	EndTime   *time.Time `json:"endTime" validate:"omitempty"`
}
