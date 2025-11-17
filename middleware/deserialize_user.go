package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"todo_list_api/api"
	"todo_list_api/utils"

	"github.com/gin-gonic/gin"
)

func DeserializeUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		authorizationHeader := c.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		}

		if token == "" {
			api.ResponseJSON(c, http.StatusUnauthorized, "Not logged in", nil)
			return
		}

		sub, err := utils.ValidateToken(token, os.Getenv("TOKEN_SECRET"))
		if err != nil {
			api.ResponseJSON(c, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		var user api.User
		result := api.DB.First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			api.ResponseJSON(c, http.StatusForbidden, "The user belonging to this token no longer exists", nil)
			return
		}

		c.Next()
	}
}
