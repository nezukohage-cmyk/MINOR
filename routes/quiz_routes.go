package routes

import (
	"lexxi/controllers"
	"lexxi/middleware"

	"github.com/gin-gonic/gin"
)

func QuizRoutes(r *gin.Engine) {
	q := r.Group("/quiz")
	q.Use(middleware.AuthRequired())

	q.POST("/start", controllers.StartQuiz)
	//q.GET("/:quiz_id/questions", controllers.QuizPage)
	q.POST("/submit", controllers.SubmitQuiz)
}
