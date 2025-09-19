package balance

import "time"

type SnapshotCreatedDomainEvent struct {
	UserID       int64
	SnapshotDate time.Time
}
