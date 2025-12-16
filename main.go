package main

import (
	"fmt"
	//"lexxi/config"
	database "lexxi/database/migrations"
	"lexxi/services"

	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	//"lexxi/middleware"
	"lexxi/routes"
)

// func main() {
// 	LoadEnv()

// 	database.Connect()
// 	services.InitCloudinary()

// 	fmt.Println("Backend oodutta eede")

// 	r := gin.Default()

// 	r.Use(func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	})
// 	r.Use(cors.New(cors.Config{
// 		AllowOrigins:     []string{"*"},
// 		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
// 		AllowCredentials: true,
// 	}))

// 	// Routes
// 	r.Use(middleware.CORSMiddleware())
// 	routes.AuthRoutes(r)
// 	routes.NoteRoutes(r)
// 	routes.UserRoutes(r)
// 	routes.SearchRoutes(r)
// 	routes.MyNotesRoutes(r)
// 	routes.CommunityRoutes(r)
// 	routes.CommentRoutes(r)
// 	routes.TagRoutes(r)
// 	routes.QuestionRoutes(r)
// 	routes.QuizRoutes(r)
// 	routes.FeedRoutes(r)
// 	routes.ChatRoutes(r)
// 	routes.PlannerRoutes(r)
// 	routes.TodoRoutes(r)
// 	// body, err := services.ListGeminiModels()
// 	// fmt.Println("List models err:", err)
// 	// fmt.Println("List models output:", body)
// 	r.Run(":8080")
// }

func main() {
	LoadEnv()
	database.Connect()
	services.InitCloudinary()
	fmt.Println("DEBUG API KEY:", os.Getenv("GEMINI_API_KEY"))

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// Routes
	routes.AuthRoutes(r)
	//routes.NoteRoutes(r)
	routes.UserRoutes(r)
	routes.SearchRoutes(r)
	routes.MyNotesRoutes(r)
	routes.CommunityRoutes(r)
	routes.CommentRoutes(r)
	routes.TagRoutes(r)
	//routes.QuestionRoutes(r)
	routes.QuizRoutes(r)
	routes.FeedRoutes(r)
	routes.ChatRoutes(r)
	routes.PlannerRoutes(r)
	routes.TodoRoutes(r)
	routes.ClusterRoutes(r)

	r.Run(":8080")
}
