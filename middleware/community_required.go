package middleware

import (
	"net/http"
	//"strings"

	"lexxi/models"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
)

func RequireCommunityMembership() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		user := &models.User{}
		if err := mgm.Coll(user).FindByID(userID, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
			c.Abort()
			return
		}

		if len(user.CommunityIDs) == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "join a community before posting"})
			c.Abort()
			return
		}

		c.Next()
	}
}
