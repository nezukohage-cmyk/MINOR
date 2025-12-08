package controllers

import (
	"net/http"

	"lexxi/services"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // username/email/phone
	Password   string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := services.Login(req.Identifier, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
	})
}
