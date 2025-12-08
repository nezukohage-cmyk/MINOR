package controllers

import (
	"net/http"

	"lexxi/services"

	"github.com/gin-gonic/gin"
)

type SignupRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required"`
}

func Signup(c *gin.Context) {
	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := services.RegisterUser(req.Username, req.Email, req.Phone, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "account created successfully",
		"token":   token,
	})
}
