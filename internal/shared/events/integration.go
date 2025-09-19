package events

import "time"

type UserProcessIntegrationEvent struct {
	UserID int64 `json:"user_id"`
}

type SnapshotCreatedIntegrationEvent struct {
	UserID       int64     `json:"user_id"`
	SnapshotDate time.Time `json:"snapshot_date"`
}
