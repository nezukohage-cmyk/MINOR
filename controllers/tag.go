package controllers

import (
	//"fmt"
	"lexxi/models"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func SeedTags(c *gin.Context) {
	tags := []models.Tag{
		{Name: "DSA", Type: "subject"},
		{Name: "OS", Type: "subject"},
		{Name: "DBMS", Type: "subject"},
		{Name: "Computer Networks", Type: "subject"},

		{Name: "Linked List", Type: "topic"},
		{Name: "Trees", Type: "topic"},
		{Name: "Paging", Type: "topic"},
		{Name: "Scheduling", Type: "topic"},
	}

	for _, t := range tags {
		var existing models.Tag
		err := mgm.Coll(&models.Tag{}).First(bson.M{"name": t.Name}, &existing)

		if err == nil {
			continue
		}

		_ = mgm.Coll(&models.Tag{}).Create(&t)
	}

	c.JSON(200, gin.H{
		"message": "tags seeded successfully",
	})
}
