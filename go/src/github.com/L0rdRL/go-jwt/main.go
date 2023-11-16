package main

import (
	"github.com/controllers"
	"github.com/gin-gonic/gin"
	"github.com/initializers"
	"github.com/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDataBase()
}
func main() {
	r := gin.Default()

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.POST("/documents", middleware.RequireAuth, controllers.AddDocument)
	r.GET("/documents", middleware.RequireAuth, controllers.GetDocuments)

	r.Run()
}
