package routes

import (
	"lexxi/controllers"

	"github.com/gin-gonic/gin"
)

func TagRoutes(r *gin.Engine) {
	tag := r.Group("/tags")

	tag.GET("/", controllers.GetTags)
	tag.POST("/create", controllers.CreateTag)
	tag.POST("/create-topics", controllers.CreateTopicsBulk)

	// OPTIONAL: remove if not needed
	// tag.GET("/seed", controllers.SeedTags)
}
