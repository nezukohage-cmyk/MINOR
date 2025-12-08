package routes

import (
	"lexxi/controllers"
	"lexxi/middleware"

	"github.com/gin-gonic/gin"
)

func FeedRoutes(r *gin.Engine) {
	feed := r.Group("/feed")
	feed.GET("/home", middleware.AuthRequired(), controllers.HomeFeed)
}
