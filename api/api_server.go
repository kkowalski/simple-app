package api

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"simple/distance"
)

func SetupServer(db *mongo.Database) *gin.Engine {
	r := gin.Default()

	distanceStore := distance.NewDistanceStore(db)
	DistanceEndpoints(r, distanceStore)

	return r
}
