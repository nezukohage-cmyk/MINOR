package routes

import (
	"lexxi/controllers"
	"lexxi/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	auth.POST("/signup", controllers.Signup)
	auth.POST("/login", controllers.Login)
	auth.GET("/me", middleware.AuthRequired(), controllers.Me)
}
