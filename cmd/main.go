package main

import (
	"log"
	"todo_list_api/api"
	"todo_list_api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to retrieve environment variables:", err)
	}

	api.InitDB()
	router := gin.Default()

	// Routes
	router.POST("/register", api.RegisterUser)
	router.POST("/login", api.LoginUser)
	router.POST("/todos", middleware.DeserializeUser(), api.CreateTodo)
	router.GET("/todos", middleware.DeserializeUser(), api.GetTodos)
	router.PUT("/todos/:id", middleware.DeserializeUser(), api.UpdateTodo)
	router.DELETE("/todos/:id", middleware.DeserializeUser(), api.DeleteTodo)

	router.Run(":8080")
}
