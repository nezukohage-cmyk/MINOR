package controllers

import (
	"lexxi/models"

	"github.com/kamva/mgm/v3"
	//"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

type InterestsRequest struct {
	Interests []string `json:"interests" binding:"required"`
}

func SaveInterests(c *gin.Context) {
	userID := c.GetString("user_id")

	var req InterestsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// find user
	var user models.User
	err := mgm.Coll(&user).FindByID(userID, &user)
	if err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	user.Interests = req.Interests

	err = mgm.Coll(&user).Update(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to save"})
		return
	}

	c.JSON(200, gin.H{"message": "interests saved"})
}
