package acceptance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"simple/api"
	"simple/distance"
	"simple/internal/testing/mongofixture"
)

func Test_Distance_API_Works(t *testing.T) {
	ctx := context.Background()

	db, cleanup := mongofixture.SetupMongo(ctx, t)
	defer cleanup()

	r := gin.Default()
	api.DistanceEndpoints(r, distance.NewDistanceStore(db))
	server := httptest.NewServer(r)
	defer server.Close()

	request := &api.DistanceRequest{UserID: "123", Timestamp: 123.123, Distance: 100}
	requestBody, err := json.Marshal(request)
	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/distance", server.URL), bytes.NewReader(requestBody))
	assert.Nil(t, err)
	response, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, response.StatusCode, 200)

	var savedDistance distance.Distance
	err = db.Collection("distance").FindOne(ctx, bson.D{}).Decode(&savedDistance)
	assert.Nil(t, err)
	assert.Equal(t, request.UserID, savedDistance.UserID)
	assert.Equal(t, request.Timestamp, savedDistance.Timestamp)
	assert.Equal(t, request.Distance, savedDistance.Distance)
}
