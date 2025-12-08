package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"

	"lexxi/models"
)

func ToggleSave(c *gin.Context) {
	userID := c.GetString("user_id")
	noteID := c.Param("id")

	user := &models.User{}
	if err := mgm.Coll(user).FindByID(userID, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}

	for i, saved := range user.SavedNoteIDs {
		if saved == noteID {

			user.SavedNoteIDs = append(user.SavedNoteIDs[:i], user.SavedNoteIDs[i+1:]...)
			mgm.Coll(user).Update(user)

			c.JSON(http.StatusOK, gin.H{
				"saved": false,
			})
			return
		}
	}
	user.SavedNoteIDs = append(user.SavedNoteIDs, noteID)
	if err := mgm.Coll(user).Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"saved": true,
	})
}
func GetSavedNotes(c *gin.Context) {
	userID := c.GetString("user_id")

	user := &models.User{}
	if err := mgm.Coll(user).FindByID(userID, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(user.SavedNoteIDs),
		"notes": user.SavedNoteIDs,
	})
}
