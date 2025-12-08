package routes

import (
	"lexxi/controllers"
	"lexxi/middleware"
	"lexxi/services"
	"os"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(r *gin.Engine) {

	chat := r.Group("/chat")

	// All chat actions require login
	chat.Use(middleware.AuthRequired())

	// ---------- AI Chat ----------
	chat.POST("/ask", controllers.AskChatbot)

	chat.GET("/env-test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"GEMINI_API_KEY": os.Getenv("GEMINI_API_KEY"),
		})
	})

	chat.GET("/models", func(c *gin.Context) {
		out, err := services.ListGeminiModels()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"models": out})
	})

	// ---------- Chat Sessions ----------
	chat.POST("/session/new", controllers.NewChatSession)
	chat.GET("/session/list", controllers.ListChatSessions)
	chat.GET("/session/:id", controllers.GetChatSession)
	chat.POST("/session/:id/send", controllers.SendChatSessionMessage)
	chat.DELETE("/session/:id", controllers.DeleteChatSession)
}
