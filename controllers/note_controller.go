package controllers

import (
	//"fmt"
	"net/http"
	"os"
	"path/filepath"

	"lexxi/models"
	"lexxi/services"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
)

func UploadNote(c *gin.Context) {
	userID := c.GetString("user_id")
	user := &models.User{}
	if err := mgm.Coll(user).FindByID(userID, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}

	// must be in a community
	if len(user.CommunityIDs) == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Join a community first"})
		return
	}

	// form fields
	subjectIDs := c.PostFormArray("subjects")
	topicIDs := c.PostFormArray("topics")

	files, err := c.FormFile("files")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file missing"})
		return
	}

	// Save temp file
	tempPath := filepath.Join(os.TempDir(), files.Filename)
	if err := c.SaveUploadedFile(files, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "temp save failed"})
		return
	}

	// UPLOAD TO CLOUDINARY
	uploadedUrl, publicID, err := services.UploadFile(tempPath, subjectIDs, topicIDs, userID)
	os.Remove(tempPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save DB entry EXACTLY like old format
	note := &models.Note{
		UserID:           userID,
		SubjectIDs:       subjectIDs,
		TopicIDs:         topicIDs,
		FileName:         files.Filename,
		FileType:         files.Header.Get("Content-Type"),
		Url:              uploadedUrl,
		CloudinaryID:     publicID,
		Size:             files.Size,
		ModerationStatus: "approved",
		CommunityID:      user.CommunityIDs[0],
		Score:            0,
		Upvotes:          []string{},
		Downvotes:        []string{},
	}

	if err := mgm.Coll(note).Create(note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save metadata"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "upload ok",
		"url":     uploadedUrl,
		"id":      note.ID.Hex(),
	})
}
