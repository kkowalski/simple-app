package mongofixture

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func SetupMongo(ctx context.Context, t *testing.T) (db *mongo.Database, cleanup func()) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctxWithTimeout, options.Client().
		ApplyURI("mongodb://localhost:27017").
		SetAuth(options.Credential{
			AuthSource: "admin",
			Password:   "simple",
			Username:   "simple"}))

	if err != nil {
		panic(err)
	}

	err = client.Ping(ctxWithTimeout, readpref.Primary())

	if err != nil {
		panic(err)
	}

	dbName := fmt.Sprintf("simple-%s", t.Name())
	cleanup = func() {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err := client.Database(dbName).Drop(ctxWithTimeout)
		if err != nil {
			panic(err)
		}
		if err = client.Disconnect(ctxWithTimeout); err != nil {
			panic(err)
		}
	}

	logrus.Info(dbName)
	db = client.Database(dbName)
	return db, cleanup
}
