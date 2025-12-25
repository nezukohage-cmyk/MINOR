package routes

import (
	"lexxi/controllers"
	"lexxi/middleware"

	"github.com/gin-gonic/gin"
)

func NotesRoutes(r *gin.Engine) {
	n := r.Group("/notes")
	n.Use(middleware.AuthRequired())

	n.POST("/save", controllers.SaveNote)
	n.GET("/saved", controllers.GetSavedNotes)
	n.POST("/unsave", controllers.UnsaveNote)
}
