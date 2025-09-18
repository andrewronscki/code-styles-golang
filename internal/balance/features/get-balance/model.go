package getbalance

import "time"

type Model struct {
	UserID       int64     `json:"user_id"`
	Balance      float64   `json:"balance"`
	SnapshotDate time.Time `json:"snapshot_date"`
}
