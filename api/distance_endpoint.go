package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"simple/distance"
)

type DistanceRequest struct {
	UserID    string  `json:"userId"`
	Timestamp float64 `json:"timestamp"`
	Distance  int64   `json:"distance"`
}

func (r DistanceRequest) ToDistance() *distance.Distance {
	return &distance.Distance{UserID: r.UserID, Timestamp: r.Timestamp, Distance: r.Distance}
}

func DistanceEndpoints(r *gin.Engine, store *distance.Store) {
	r.POST("/distance", func(c *gin.Context) {
		var req DistanceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx := context.Background()

		err := store.Save(ctx, req.ToDistance())
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"your_distance": req.Distance})
	})
}
