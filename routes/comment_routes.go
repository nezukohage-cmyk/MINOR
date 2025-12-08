package routes

import (
	"lexxi/controllers"
	"lexxi/middleware"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(r *gin.Engine) {
	comment := r.Group("/notes/:id/comments")
	comment.Use(middleware.AuthRequired())

	comment.POST("/", controllers.CreateComment)
	comment.GET("/", controllers.GetComments)
}
