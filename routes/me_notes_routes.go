package routes

import (
	"lexxi/controllers"
	"lexxi/middleware"

	"github.com/gin-gonic/gin"
)

func MyNotesRoutes(r *gin.Engine) {
	notes := r.Group("/me")
	notes.Use(middleware.AuthRequired())
	notes.GET("/notes", controllers.GetMyNotes)
}
