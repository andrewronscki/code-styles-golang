package balance

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ba *Balance) MarshalBSON() ([]byte, error) {
	b := struct {
		ID           primitive.ObjectID `bson:"_id,omitempty"`
		UserID       int64              `bson:"user_id"`
		Balance      float64            `bson:"balance"`
		SnapshotDate time.Time          `bson:"snapshot_date"`
	}{
		UserID:       ba.UserID,
		Balance:      ba.Balance,
		SnapshotDate: ba.SnapshotDate,
	}

	if ba.ID != "" {
		if id, err := primitive.ObjectIDFromHex(ba.ID); err == nil {
			b.ID = id
		}
	}

	return bson.Marshal(b)
}

func (ba *Balance) UnmarshalBSON(data []byte) error {
	var b struct {
		ID           string    `bson:"_id,omitempty"`
		UserID       int64     `bson:"user_id"`
		Balance      float64   `bson:"balance"`
		SnapshotDate time.Time `bson:"snapshot_date"`
	}

	if err := bson.Unmarshal(data, &b); err != nil {
		return err
	}

	ba.ID = b.ID
	ba.UserID = b.UserID
	ba.Balance = b.Balance
	ba.SnapshotDate = b.SnapshotDate

	return nil
}
