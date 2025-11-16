package api

import (
	"log"
	"net/http"
	"os"
	"todo_list_api/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	if err := DB.AutoMigrate(&Todo{}); err != nil {
		log.Fatal("Failed to migrate schema:", err)
	}
}

func CreateUser(c *gin.Context) {
	var user User

	// Bind the request body
	if err := c.ShouldBindJSON(&user); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	// Hash the password first
	user.Password, _ = utils.HashPassword(user.Password)

	DB.Create(&user)
	// TODO: Create JWT token if registration is successful
	ResponseJSON(c, http.StatusCreated, "Todo created successfully", user)
}

func CreateTodo(c *gin.Context) {
	var todo Todo

	// Bind the request body
	if err := c.ShouldBindJSON(&todo); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}
	DB.Create(&todo)
	ResponseJSON(c, http.StatusCreated, "Todo created successfully", todo)
}

// TODO: Implement pagination and filtering
func GetTodos(c *gin.Context) {
	var todos []Todo
	DB.Find(&todos)
	ResponseJSON(c, http.StatusOK, "Todos retrieved successfully", todos)
}

func UpdateTodo(c *gin.Context) {
	var todo Todo
	if err := DB.First(&todo, c.Param("id")).Error; err != nil {
		ResponseJSON(c, http.StatusNotFound, "Todo not found", nil)
		return
	}

	// Bind the request body
	if err := c.ShouldBindJSON(&todo); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	DB.Save(&todo)
	ResponseJSON(c, http.StatusOK, "Todo updated successfully", todo)
}

func DeleteTodo(c *gin.Context) {
	var todo Todo
	if err := DB.Delete(&todo, c.Param("id")).Error; err != nil {
		ResponseJSON(c, http.StatusNotFound, "Todo not found", nil)
		return
	}
	ResponseJSON(c, http.StatusOK, "Todo deleted successfully", nil)
}
