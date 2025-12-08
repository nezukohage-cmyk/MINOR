package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Me(c *gin.Context) {
	userId := c.GetString("user_id")

	c.JSON(http.StatusOK, gin.H{
		"message": "authenticated request",
		"user_id": userId,
	})
}
