// main.go
package main

import (
	"fmt"
	"log"
	"notes_application/controllers"
	"notes_application/initializers"
	"notes_application/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.InitialMigrations()
}

func main() {
	fmt.Println("Notes API")

	r := gin.Default()

	// Routes
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)

	// Protected routes using AuthMiddleware
	protected := r.Group("/notes")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("", controllers.GetAllNotes)
		protected.POST("", controllers.CreateNote)
		protected.DELETE("", controllers.DeleteNote)
	}

	// Run the server
	if err := r.Run(); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
