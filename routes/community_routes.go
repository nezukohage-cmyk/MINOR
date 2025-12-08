package routes

import (
	"lexxi/controllers"
	"lexxi/middleware"

	"github.com/gin-gonic/gin"
)

func CommunityRoutes(r *gin.Engine) {
	group := r.Group("/communities")
	group.Use(middleware.AuthRequired())
	group.POST("/create", controllers.CreateCommunity)
	group.GET("/search", controllers.SearchCommunities)
	group.POST("/join/:name", controllers.JoinCommunity)
	group.POST("/leave/:name", controllers.LeaveCommunity)

	group.GET("/me", controllers.MyCommunities)
}
