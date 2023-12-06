package middleware

import (
	"net/http"
	"restful-portal/src/helpers"
	"restful-portal/src/modules/users"
	"strings"

	helperJWT "restful-portal/src/helpers/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(authService helperJWT.Service, userService users.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helpers.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""

		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateAccessToken(tokenString)

		if err != nil {
			response := helpers.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helpers.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)

		if err != nil {
			response := helpers.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}
