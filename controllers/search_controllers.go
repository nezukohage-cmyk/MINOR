// package controllers

// import (
// 	"net/http"

// 	"lexxi/models"

// 	"github.com/gin-gonic/gin"
// 	"github.com/kamva/mgm/v3"
// 	"go.mongodb.org/mongo-driver/bson"
// )

// func SearchNotes(c *gin.Context) {

// 	query := c.Query("q")

// 	if query == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "search query 'q' required"})
// 		return
// 	}

// 	// MongoDB OR search across multiple fields
// 	filter := bson.M{
// 		"$or": []bson.M{
// 			{"file_name": bson.M{"$regex": query, "$options": "i"}},
// 			{"file_type": bson.M{"$regex": query, "$options": "i"}},
// 			{"subject_id": bson.M{"$regex": query, "$options": "i"}},
// 			{"topic_id": bson.M{"$regex": query, "$options": "i"}},
// 		},
// 	}

// 	var results []models.Note

// 	err := mgm.Coll(&models.Note{}).SimpleFind(&results, filter)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"count":   len(results),
// 		"results": results,
// 	})
// }

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

func SearchNotes(c *gin.Context) {
	query := c.Query("q")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search query 'q' required"})
		return
	}

	// Pagination params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := (page - 1) * limit

	// MongoDB filter
	filter := bson.M{
		"$or": []bson.M{
			{"file_name": bson.M{"$regex": query, "$options": "i"}},
			{"subject_id": bson.M{"$regex": query, "$options": "i"}},
			{"topic_id": bson.M{"$regex": query, "$options": "i"}},
			{"file_type": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	// Query with pagination
	opts := options.Find()
	opts.SetSkip(int64(skip))
	opts.SetLimit(int64(limit))

	var results []models.Note

	err := mgm.Coll(&models.Note{}).SimpleFind(&results, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	// Count total matching docs for frontend pagination UI
	total, _ := mgm.Coll(&models.Note{}).CountDocuments(mgm.Ctx(), filter)

	c.JSON(http.StatusOK, gin.H{
		"page":    page,
		"limit":   limit,
		"total":   total,
		"pages":   (total + int64(limit) - 1) / int64(limit),
		"results": results,
	})
}
