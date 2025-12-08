package routes

import (
	"lexxi/controllers"
	"lexxi/middleware"

	"github.com/gin-gonic/gin"
)

func QuestionRoutes(r *gin.Engine) {
	q := r.Group("/questions")
	q.Use(middleware.AuthRequired()) // login required

	q.POST("/create", controllers.CreateQuestion)
}
