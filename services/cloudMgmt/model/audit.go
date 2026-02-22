package model

import "time"

type AuditTrailModel struct {
	UserID         int64
	ChronoSequence string
	Action         string
	ServerID       string
	Path           string
	IPAddress      string
	ResStatus      int
	OldValue       string
	NewValue       string
	ActionTime     time.Time
}
