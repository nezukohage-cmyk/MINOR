package routes

import (
	"lexxi/controllers"
	"lexxi/middleware"

	"github.com/gin-gonic/gin"
)

func SearchRoutes(r *gin.Engine) {
	search := r.Group("/search")
	search.GET("/", middleware.AuthRequired(), controllers.SearchNotes)
}
