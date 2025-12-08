package main

import (
	"fmt"
	//"lexxi/config"
	database "lexxi/database/migrations"
	"lexxi/services"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"lexxi/routes"
)

func main() {
	LoadEnv()
	fmt.Println("DEBUG API KEY:", os.Getenv("GEMINI_API_KEY"))

	database.Connect()
	services.InitCloudinary()

	fmt.Println("Backend oodutta eede")

	r := gin.Default()

	// -----------------------------
	// 🔥 Enable CORS for Flutter Web
	// -----------------------------
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Routes
	routes.AuthRoutes(r)
	routes.NoteRoutes(r)
	routes.UserRoutes(r)
	routes.SearchRoutes(r)
	routes.MyNotesRoutes(r)
	routes.CommunityRoutes(r)
	routes.CommentRoutes(r)
	routes.TagRoutes(r)
	routes.QuestionRoutes(r)
	routes.QuizRoutes(r)
	routes.FeedRoutes(r)
	routes.ChatRoutes(r)
	routes.PlannerRoutes(r)
	routes.TodoRoutes(r)
	// body, err := services.ListGeminiModels()
	// fmt.Println("List models err:", err)
	// fmt.Println("List models output:", body)
	r.Run(":8080")
}
