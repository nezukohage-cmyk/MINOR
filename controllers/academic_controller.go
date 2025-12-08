package controllers

import (
	"log"
	"net/http"

	"lexxi/models"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func Subjects(c *gin.Context) {
	var subs []models.Subject
	err := mgm.Coll(&models.Subject{}).SimpleFind(&subs, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, subs)
}

func Topics(c *gin.Context) {
	subject := c.Query("subject")

	filter := bson.M{}
	if subject != "" {
		filter["subject"] = subject
	}

	var topics []models.Topic
	err := mgm.Coll(&models.Topic{}).SimpleFind(&topics, filter)
	if err != nil {
		log.Println("topics fetch err:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, topics)
}
