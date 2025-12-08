package controllers

import (
	"net/http"
	"strconv"

	"lexxi/models"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMyNotes(c *gin.Context) {
	userID := c.GetString("user_id")

	// pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	filter := bson.M{
		"user_id": userID,
	}

	opts := options.Find()
	opts.SetSkip(int64(skip))
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.M{"created_at": -1}) // newest first

	var notes []models.Note

	err := mgm.Coll(&models.Note{}).SimpleFind(&notes, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed fetching notes"})
		return
	}

	total, _ := mgm.Coll(&models.Note{}).CountDocuments(mgm.Ctx(), filter)

	c.JSON(http.StatusOK, gin.H{
		"page":    page,
		"limit":   limit,
		"total":   total,
		"results": notes,
	})
}
