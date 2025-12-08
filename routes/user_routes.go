package routes

import (
	"lexxi/controllers"
	"lexxi/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	user := router.Group("/user")
	user.Use(middleware.AuthRequired())
	{
		user.POST("/interests", controllers.SaveInterests)
	}
}
