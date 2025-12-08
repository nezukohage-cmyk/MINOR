package routes

import (
	"lexxi/controllers"
	"lexxi/middleware"

	//"lexxi/services"
	"github.com/gin-gonic/gin"
)

func NoteRoutes(r *gin.Engine) {
	notes := r.Group("/notes")
	notes.POST("/upload", middleware.AuthRequired(), controllers.UploadNote)
	notes.POST("/:id/upvote", middleware.AuthRequired(), controllers.Upvote)
	notes.POST("/:id/downvote", middleware.AuthRequired(), controllers.Downvote)
	notes.POST("/:id/save", middleware.AuthRequired(), controllers.ToggleSave)
	notes.GET("/saved", middleware.AuthRequired(), controllers.GetSavedNotes)
	//notes.POST("/test-moderation", middleware.AuthRequired(), services.ModerateImageHF)

}
