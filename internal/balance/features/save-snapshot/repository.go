package savesnapshot

import (
	"context"

	balance "github.com/andre/code-styles-golang/internal/balance/domain"
	"github.com/andre/code-styles-golang/pkg/datadog/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	Insert(ctx context.Context, balance *balance.Balance) (id string, err error)
}

type BalanceRepositoryParams struct {
	Database   string
	Collection string
}

type BalanceRepository struct {
	database   string
	collection string
}

func (r *BalanceRepository) Insert(ctx context.Context, balance *balance.Balance) (id string, err error) {
	db := mongodb.Client().Database(r.database)

	col := db.Collection(r.collection)

	balance.MarshalBSON()

	result, err := col.InsertOne(ctx, balance)

	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func NewBalanceRepository(p *BalanceRepositoryParams) Repository {
	return &BalanceRepository{
		database:   p.Database,
		collection: p.Collection,
	}
}
