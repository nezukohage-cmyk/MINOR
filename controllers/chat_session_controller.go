package controllers

import (
	"lexxi/models"
	"lexxi/services"

	//"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

// ------------------------------
// CREATE NEW CHAT SESSION
// POST /chat/session/new
// ------------------------------
func NewChatSession(c *gin.Context) {
	userID := c.GetString("user_id")

	var body struct {
		Title string `json:"title"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	chat := &models.Chat{
		UserID:   userID,
		Title:    body.Title,
		Messages: []models.ChatMessage{},
	}

	if err := mgm.Coll(chat).Create(chat); err != nil {
		c.JSON(500, gin.H{"error": "failed to create chat session"})
		return
	}

	c.JSON(200, gin.H{"chat_id": chat.ID.Hex()})
}

// ------------------------------
// LIST USER CHAT SESSIONS
// GET /chat/session/list
// ------------------------------
func ListChatSessions(c *gin.Context) {
	userID := c.GetString("user_id")

	var chats []models.Chat
	err := mgm.Coll(&models.Chat{}).SimpleFind(&chats, bson.M{
		"user_id": userID,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch chats"})
		return
	}

	c.JSON(200, gin.H{"chats": chats})
}

// ------------------------------
// GET CHAT MESSAGES
// GET /chat/session/:id
// ------------------------------
func GetChatSession(c *gin.Context) {
	chatID := c.Param("id")

	var chat models.Chat
	err := mgm.Coll(&chat).FindByID(chatID, &chat)

	if err != nil {
		c.JSON(404, gin.H{"error": "chat not found"})
		return
	}

	c.JSON(200, gin.H{"chat": chat})
}

// ------------------------------
// SEND MESSAGE + AI RESPONSE
// POST /chat/session/:id/send
// ------------------------------
func SendChatSessionMessage(c *gin.Context) {
	chatID := c.Param("id")

	var chat models.Chat
	if err := mgm.Coll(&chat).FindByID(chatID, &chat); err != nil {
		c.JSON(404, gin.H{"error": "chat not found"})
		return
	}

	var body struct {
		Message string `json:"message"`
	}

	if err := c.BindJSON(&body); err != nil || body.Message == "" {
		c.JSON(400, gin.H{"error": "message is required"})
		return
	}

	// USER MESSAGE STORED
	chat.Messages = append(chat.Messages, models.ChatMessage{
		Role: "user",
		Text: body.Message,
	})

	// AI REPLY
	reply, err := services.AskGemini("", body.Message)
	if err != nil {
		c.JSON(500, gin.H{"error": "AI failed to respond"})
		return
	}

	// STORE AI MESSAGE
	chat.Messages = append(chat.Messages, models.ChatMessage{
		Role: "assistant",
		Text: reply,
	})

	mgm.Coll(&chat).Update(&chat)

	c.JSON(200, gin.H{
		"reply": reply,
	})
}

// ------------------------------
// DELETE CHAT SESSION
// DELETE /chat/session/:id
// ------------------------------
func DeleteChatSession(c *gin.Context) {
	chatID := c.Param("id")

	var chat models.Chat
	if err := mgm.Coll(&chat).FindByID(chatID, &chat); err != nil {
		c.JSON(404, gin.H{"error": "chat not found"})
		return
	}

	mgm.Coll(&chat).Delete(&chat)

	c.JSON(200, gin.H{"message": "chat deleted"})
}
