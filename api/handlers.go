package api

import (
	"log"
	"net/http"
	"os"
	"strings"
	"todo_list_api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DB_URL")
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	if err := DB.AutoMigrate(&Todo{}); err != nil {
		log.Fatal("Failed to migrate schema:", err)
	}
}

func RegisterUser(c *gin.Context) {
	var user User

	// Bind the request body
	if err := c.ShouldBindJSON(&user); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	// Hash the password first
	user.Password, _ = utils.HashPassword(user.Password)

	DB.Create(&user)

	token, _ := utils.GenerateToken(user.ID, os.Getenv("TOKEN_SECRET"))
	ResponseJSON(c, http.StatusCreated, "User registration successful", token)
}

func LoginUser(c *gin.Context) {
	var loginInput LoginInput

	// Bind the request body
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	var user User

	// Check if email exists in database
	result := DB.First(&user, "email = ?", strings.ToLower(loginInput.Email))
	if result.Error != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid email or password", nil)
		return
	}

	// Check if password matches
	if !utils.VerifyPassword(user.Password, loginInput.Password) {
		ResponseJSON(c, http.StatusBadRequest, "Invalid email or password", nil)
		return
	}

	token, _ := utils.GenerateToken(loginInput.Email, os.Getenv("TOKEN_SECRET"))
	ResponseJSON(c, http.StatusOK, "User login successful", token)
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
