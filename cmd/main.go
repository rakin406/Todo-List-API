package main

import (
	"todo_list_api/api"

	"github.com/gin-gonic/gin"
)

func main() {
	api.InitDB()
	router := gin.Default()

	// Routes
	router.POST("/todos", api.CreateTodo)
	router.GET("/todos", api.GetTodos)
	router.PUT("/todos/:id", api.UpdateTodo)
	router.DELETE("/todos/:id", api.DeleteTodo)

	router.Run(":8080")
}
