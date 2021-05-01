package db

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Setup(ctx context.Context) (db *mongo.Database, cleanup func()) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctxWithTimeout, options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetAuth(options.Credential{
			AuthSource: "admin",
			Password:   "simple",
			Username:   "simple"}))
	if err != nil {
		panic(fmt.Errorf("error connecting to mongo: %w", err))
	}

	err = client.Ping(ctxWithTimeout, readpref.Primary())
	if err != nil {
		panic(fmt.Errorf("error pinging mongo: %w", err))
	}

	cleanup = func() {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err = client.Disconnect(ctxWithTimeout); err != nil {
			logrus.Error("error disconnecting: %w", err)
		}
	}

	db = client.Database("simple")
	return db, cleanup
}
