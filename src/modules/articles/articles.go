package articles

import (
	"fmt"
	"net/http"
	"restful-portal/src/helpers"
	"restful-portal/src/modules/users"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) GetArticles(c *gin.Context) {
	AuthorID, _ := strconv.Atoi(c.Query("author_id"))

	articles, err := h.service.GetAll(AuthorID)
	if err != nil {
		response := helpers.APIResponse("Error to get articles", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helpers.APIResponse("List of articles", http.StatusOK, "success", FormatArticles(articles))
	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetArticle(c *gin.Context) {
	var input GetArticleInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helpers.APIResponse("Failed to get detail of article", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	articleDetail, err := h.service.Get(input)
	if err != nil {
		response := helpers.APIResponse("Failed to get detail of article", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if articleDetail.ID == 0 {
		response := helpers.APIResponse("Article not found", http.StatusNotFound, "error", nil)
		c.JSON(http.StatusNotFound, response)
		return

	}

	response := helpers.APIResponse("Article Detail", http.StatusOK, "success", FormatArticle(articleDetail))
	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateArticle(c *gin.Context) {
	var input CreateArticleInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helpers.FormatValidationError(err)
		errorMessage := gin.H{"error": errors}

		response := helpers.APIResponse("Failed to create article", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(users.User)
	input.User = currentUser

	slugArticle := slug.Make(fmt.Sprintf("%s %d", input.Title, input.User.ID))
	isSlugExist, err := h.service.IsSlugAvailable(slugArticle)
	if err != nil || !isSlugExist {
		response := helpers.APIResponse("Failed to create article because already exist", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	new, err := h.service.Create(input, slugArticle)
	if err != nil {
		response := helpers.APIResponse("Failed to create article", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helpers.APIResponse("Success to create article", http.StatusOK, "success", FormatArticle(new))
	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateArticle(c *gin.Context) {
	var inputID GetArticleInput
	var inputData UpdateArticleInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helpers.APIResponse("Failed to update article", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helpers.FormatValidationError(err)
		errorMessage := gin.H{"error": errors}

		response := helpers.APIResponse("Failed to update article", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(users.User)
	inputData.User = currentUser

	deleteArticle, err := h.service.Update(inputID, inputData)
	if err != nil {
		response := helpers.APIResponse("Failed to update article", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helpers.APIResponse("Success to update article", http.StatusOK, "success", FormatArticle(deleteArticle))
	c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteArticle(c *gin.Context) {
	var inputID GetArticleInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helpers.APIResponse("Failed to delete article", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(users.User)

	_, err = h.service.Delete(inputID, currentUser)
	if err != nil {
		response := helpers.APIResponse("Failed to update article", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helpers.APIResponse("Success to update article", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
