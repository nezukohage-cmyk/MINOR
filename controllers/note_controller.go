// package controllers

// import (
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"path/filepath"

// 	//"strconv"

// 	"github.com/gin-gonic/gin"
// 	//"lexxi/middleware"
// 	"lexxi/models"
// 	"lexxi/services"

// 	"github.com/kamva/mgm/v3"
// )

// approved, reason := services.ModerateFile(tempFilePath)
// if !approved {
//     // delete the temporary file
//     os.Remove(tempFilePath)

//     c.JSON(400, gin.H{
//         "success": false,
//         "message": "Moderation failed",
//         "reason":  reason,
//     })
//     return
// }

// func UploadNote(c *gin.Context) {

// 	userID := c.GetString("user_id")
// 	user := &models.User{}
// 	if err := mgm.Coll(user).FindByID(userID, user); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
// 		return
// 	}
// 	if len(user.CommunityIDs) == 0 {
// 		c.JSON(http.StatusForbidden, gin.H{
// 			"error": "You must join a community before uploading.",
// 		})
// 		return
// 	}
// 	subject := c.PostForm("subject")
// 	topic := c.PostForm("topic")

// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "file missing"})
// 		return
// 	}

// 	tempPath := filepath.Join(os.TempDir(), file.Filename)
// 	err = c.SaveUploadedFile(file, tempPath)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save temporary file"})
// 		return
// 	}
// 	uploadedUrl, publicID, err := services.UploadFile(tempPath, subject, topic, userID)
// 	if err != nil {
// 		fmt.Println(" Cloudinary error:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	note := &models.Note{
// 		UserID:           userID,
// 		SubjectIDs:       c.PostFormArray("subjects"),
// 		TopicIDs:         c.PostFormArray("topics"),
// 		FileName:         file.Filename,
// 		FileType:         file.Header.Get("Content-Type"),
// 		Size:             file.Size,
// 		Url:              uploadedUrl,
// 		CloudinaryID:     publicID,
// 		ModerationStatus: "pending",
// 	}

// 	err = mgm.Coll(note).Create(note)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save metadata"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "uploaded and pending moderation",
// 		"url":     uploadedUrl,
// 		"id":      note.ID.Hex(),
// 	})
// }

package controllers

import (
	"fmt"
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

	// User must be in a community
	if len(user.CommunityIDs) == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You must join a community before uploading.",
		})
		return
	}

	// Form data
	subjectIDs := c.PostFormArray("subjects")
	topicIDs := c.PostFormArray("topics")

	// Read file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file missing"})
		return
	}

	// Save temp file
	tempPath := filepath.Join(os.TempDir(), file.Filename)
	err = c.SaveUploadedFile(file, tempPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save temporary file"})
		return
	}

	// -------------------------------------------
	// 🔥 FILENAME-BASED STUDY FILTER
	// -------------------------------------------

	// Build keyword list from subject + topic IDs
	var keywords []string
	keywords = append(keywords, subjectIDs...)
	keywords = append(keywords, topicIDs...)

	// Validate filename against keywords
	if !services.IsFilenameStudyRelated(file.Filename, keywords) {
		os.Remove(tempPath)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"reason":  "File name does not match any study topic or subject",
		})
		return
	}
	// -------------------------------------------

	// Upload to Cloudinary
	uploadedUrl, publicID, err := services.UploadFile(tempPath, subjectIDs, topicIDs, userID)
	if err != nil {
		fmt.Println("Cloudinary error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete temp file after Cloudinary upload
	os.Remove(tempPath)

	// Create MongoDB note entry
	note := &models.Note{
		UserID:           userID,
		SubjectIDs:       subjectIDs,
		TopicIDs:         topicIDs,
		FileName:         file.Filename,
		FileType:         file.Header.Get("Content-Type"),
		Size:             file.Size,
		Url:              uploadedUrl,
		CloudinaryID:     publicID,
		ModerationStatus: "approved",
		CommunityID:      user.CommunityIDs[0],
	}

	err = mgm.Coll(note).Create(note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save metadata"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Upload successful",
		"url":     uploadedUrl,
		"id":      note.ID.Hex(),
	})
}
