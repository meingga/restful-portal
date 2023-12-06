package users

import (
	"net/http"
	"restful-portal/src/helpers"
	HelperJWT "restful-portal/src/helpers/jwt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type userHandler struct {
	userService Service
	authService HelperJWT.Service
}

func NewUserHandler(userService Service, authService HelperJWT.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (u *userHandler) Register(c *gin.Context) {
	var input RegisterUserInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helpers.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helpers.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := u.userService.Register(input)

	if err != nil {
		response := helpers.APIResponse("Register account failed", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := FormatRegisterUser(newUser)
	response := helpers.APIResponse("Account has ben registered", http.StatusOK, "Success", formatter)

	c.JSON(http.StatusOK, response)
}

func (u *userHandler) Login(c *gin.Context) {
	var input LoginUserInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helpers.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helpers.APIResponse("Login account failed", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := u.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helpers.APIResponse("Login account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, refreshToken, err := u.authService.GenerateToken(loggedinUser.ID)
	if err != nil {
		response := helpers.APIResponse("Login failed", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := FormatLoginUser(loggedinUser, token, refreshToken)

	response := helpers.APIResponse("Successfuly loggedin", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (u *userHandler) RefreshToken(c *gin.Context) {
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

	validateToken, err := u.authService.ValidateRefreshToken(tokenString)

	if err != nil {
		response := helpers.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)

		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	claim, ok := validateToken.Claims.(jwt.MapClaims)

	if !ok || !validateToken.Valid {
		response := helpers.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)

		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	userID := int(claim["user_id"].(float64))

	token, _, err := u.authService.GenerateToken(userID)
	if err != nil {
		response := helpers.APIResponse("Refresh Token failed", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := FormatRefreshTokenUser(token)

	response := helpers.APIResponse("Successfuly Refresh Token", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (u *userHandler) FetchUser(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(User)
	formatter := FormatUser(currentUser)
	response := helpers.APIResponse("Successfuly fetch user data", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (u *userHandler) Update(c *gin.Context) {
	var inputID GetUserDetailInput
	var inputData ChangePasswordUserInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helpers.APIResponse("Failed to update user", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helpers.FormatValidationError(err)
		errorMessage := gin.H{"error": errors}

		response := helpers.APIResponse("Failed to update user", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(User)
	inputData.User = currentUser

	updateUser, err := u.userService.ChangePassword(inputID, inputData)
	if err != nil {
		response := helpers.APIResponse("Failed to update user", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helpers.APIResponse("Success to update user", http.StatusOK, "success", FormatUser(updateUser))
	c.JSON(http.StatusOK, response)
}
