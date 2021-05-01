package distance

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	col *mongo.Collection
}

type Distance struct {
	UserID    string `bson:"orbId"`
	Timestamp float64
	Distance  int64
}

func (s Store) Save(ctx context.Context, distance *Distance) error {
	_, err := s.col.InsertOne(ctx, distance)
	if err != nil {
		return fmt.Errorf("error saving distance:%w", err)
	}
	return nil
}

func NewDistanceStore(db *mongo.Database) *Store {
	return &Store{col: db.Collection("distance")}
}
