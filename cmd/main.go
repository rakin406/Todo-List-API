package main

import (
	"log"
	"todo_list_api/api"

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
	router.POST("/todos", api.CreateTodo)
	router.GET("/todos", api.GetTodos)
	router.PUT("/todos/:id", api.UpdateTodo)
	router.DELETE("/todos/:id", api.DeleteTodo)

	router.Run(":8080")
}
