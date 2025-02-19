package server

import (
	"net/http"
	"yorgut/internal/controller"

	_ "yorgut/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Yorgut APIs
// @version 1.0
// @description Our very own OG Social Media.

// @in header

// @name token// @license.name Apache 2.0

// @host localhost:8080
// @BasePath /api/v1

// @schemes http
func NewRouter(userController *controller.UserController) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "Welcome")
	})

	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/user")
		{
			users.POST("", userController.Create)
			users.GET("", userController.FindMany)
			users.PUT("", userController.UpdateUser)
			users.GET("/:userId", userController.FindUserById)
			users.DELETE("/:userId", userController.DeleteUser)
		}
	}

	return router
}
