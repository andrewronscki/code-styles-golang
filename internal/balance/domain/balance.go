package balance

import (
	"time"

	"github.com/andre/code-styles-golang/internal/shared/utils"
)

type Balance struct {
	ID           string    `json:"id,omitempty"`
	UserID       int64     `json:"user_id,omitempty"`
	Balance      float64   `json:"balance"`
	SnapshotDate time.Time `json:"snapshot_date"`
}

func NewBalance(userId int64, snapshotDate time.Time) *Balance {
	return &Balance{
		UserID:       userId,
		SnapshotDate: snapshotDate,
		Balance:      0.0,
	}
}

func (b *Balance) SetBalance(balance float64) {
	b.Balance = balance
}

func (b *Balance) SetID(id string) {
	b.ID = id
}

func (b *Balance) Marshal(dest any) {
	utils.DeepCopy(b, dest)
}

func (b *Balance) RaiseSnapshotCreatedDomainEvent() *SnapshotCreatedDomainEvent {
	event := &SnapshotCreatedDomainEvent{
		UserID:       b.UserID,
		SnapshotDate: b.SnapshotDate,
	}

	return event
}
