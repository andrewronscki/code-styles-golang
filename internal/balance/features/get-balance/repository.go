package getbalance

import (
	"context"
	"errors"
	"time"

	balance "github.com/andre/code-styles-golang/internal/balance/domain"
	"github.com/andre/code-styles-golang/pkg/datadog/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Filter struct {
	UserID       int64
	SnapshotDate time.Time
}

type Repository interface {
	FindBalance(ctx context.Context, f *Filter) (*balance.Balance, error)
}

type BalanceRepositoryParams struct {
	Database   string
	Collection string
}

type BalanceRepository struct {
	database   string
	collection string
}

func (r *BalanceRepository) FindBalance(ctx context.Context, f *Filter) (*balance.Balance, error) {
	collection := mongodb.Client().Database(r.database).Collection(r.collection)

	snapshot := f.SnapshotDate

	filter := bson.M{
		"user_id":       f.UserID,
		"snapshot_date": bson.M{"$lte": snapshot},
	}

	opts := options.FindOne().SetSort(bson.D{{Key: "snapshot_date", Value: -1}})

	var out balance.Balance
	err := collection.FindOne(ctx, filter, opts).Decode(&out)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func NewBalanceRepository(p *BalanceRepositoryParams) Repository {
	return &BalanceRepository{
		database:   p.Database,
		collection: p.Collection,
	}
}
