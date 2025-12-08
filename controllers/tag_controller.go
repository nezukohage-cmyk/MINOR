package controllers

import (
	"context"
	"net/http"
	"strings"

	"lexxi/models"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTags(c *gin.Context) {
	var tags []models.Tag

	err := mgm.Coll(&models.Tag{}).SimpleFind(&tags, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tags"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

func CreateTag(c *gin.Context) {

	var body struct {
		Name           string   `json:"name"`
		Type           string   `json:"type"`
		ParentSubjects []string `json:"parent_subjects"` // names, not IDs
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	name := strings.TrimSpace(body.Name)
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	lowerName := strings.ToLower(name)

	// =========================================
	// SUBJECT CREATION
	// =========================================
	if body.Type == "subject" {

		var existing models.Tag

		err := mgm.Coll(&models.Tag{}).First(bson.M{
			"type": "subject",
			"name": bson.M{"$regex": "^" + lowerName + "$", "$options": "i"},
		}, &existing)

		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "subject already exists"})
			return
		}

		newSub := &models.Tag{
			Name:        name,
			Type:        "subject",
			ChildTopics: []string{},
		}

		if err := mgm.Coll(newSub).Create(newSub); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subject"})
			return
		}

		c.JSON(200, gin.H{
			"message": "Subject created",
			"data":    newSub,
		})
		return
	}

	// =========================================
	// TOPIC CREATION
	// =========================================
	if body.Type == "topic" {

		if len(body.ParentSubjects) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "topic must have at least one parent subject name",
			})
			return
		}

		var parentIDs []string

		// Lookup subjects by name
		for _, subName := range body.ParentSubjects {

			var subject models.Tag

			err := mgm.Coll(&models.Tag{}).First(bson.M{
				"type": "subject",
				"name": bson.M{
					"$regex":   "^" + strings.ToLower(subName) + "$",
					"$options": "i",
				},
			}, &subject)

			if err != nil {
				c.JSON(400, gin.H{"error": "parent subject not found: " + subName})
				return
			}

			parentIDs = append(parentIDs, subject.ID.Hex())
		}

		// Create topic
		newTopic := &models.Tag{
			Name:           name,
			Type:           "topic",
			ParentSubjects: parentIDs,
		}

		if err := mgm.Coll(newTopic).Create(newTopic); err != nil {
			c.JSON(500, gin.H{"error": "failed to create topic"})
			return
		}

		// Link topic to subjects
		for _, id := range parentIDs {

			objID, _ := primitive.ObjectIDFromHex(id)

			mgm.Coll(&models.Tag{}).UpdateOne(
				context.Background(),
				bson.M{"_id": objID},
				bson.M{"$addToSet": bson.M{"child_topics": newTopic.ID.Hex()}},
			)
		}

		c.JSON(200, gin.H{
			"message": "Topic created",
			"data":    newTopic,
		})
		return
	}

	c.JSON(400, gin.H{"error": "type must be subject or topic"})
}
func CreateTopicsBulk(c *gin.Context) {

	var body struct {
		ParentSubject string   `json:"parent_subject"`
		Topics        []string `json:"topics"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid JSON"})
		return
	}

	if body.ParentSubject == "" {
		c.JSON(400, gin.H{"error": "parent_subject is required"})
		return
	}

	if len(body.Topics) == 0 {
		c.JSON(400, gin.H{"error": "topics list cannot be empty"})
		return
	}

	// ---------------------------
	// FIND SUBJECT BY NAME
	// ---------------------------
	var subject models.Tag
	err := mgm.Coll(&models.Tag{}).First(bson.M{
		"type": "subject",
		"name": bson.M{"$regex": "^" + strings.ToLower(body.ParentSubject) + "$", "$options": "i"},
	}, &subject)

	if err != nil {
		c.JSON(400, gin.H{"error": "parent subject not found"})
		return
	}

	subjectID := subject.ID.Hex()
	created := []models.Tag{}

	// ---------------------------
	// CREATE MULTIPLE TOPICS
	// ---------------------------
	for _, tname := range body.Topics {

		tname = strings.TrimSpace(tname)
		if tname == "" {
			continue
		}

		newTopic := &models.Tag{
			Name:           tname,
			Type:           "topic",
			ParentSubjects: []string{subjectID},
		}

		if err := mgm.Coll(newTopic).Create(newTopic); err != nil {
			continue
		}

		created = append(created, *newTopic)

		mgm.Coll(&models.Tag{}).UpdateOne(
			context.Background(),
			bson.M{"_id": subject.ID},
			bson.M{"$addToSet": bson.M{"child_topics": newTopic.ID.Hex()}},
		)
	}

	c.JSON(200, gin.H{
		"message": "Topics created successfully",
		"count":   len(created),
		"topics":  created,
	})
}
