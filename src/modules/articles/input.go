package articles

import (
	"restful-portal/src/modules/users"
)

type GetArticleInput struct {
	ID int `uri:"id" binding:"required"`
}

type CreateArticleInput struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
	PublishedAt string `json:"published_at" binding:"required"`
	User        users.User
}

type UpdateArticleInput struct {
	Content     string `json:"content" binding:"required"`
	PublishedAt string `json:"published_at" binding:"required"`
	User        users.User
}

type CreateCommentInput struct {
	Content string `json:"content" binding:"required"`
	Article Article
	User    users.User
}
