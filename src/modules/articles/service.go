package articles

import (
	"errors"
	"restful-portal/src/modules/users"
	"time"
)

type Service interface {
	GetAll(userID int) ([]Article, error)
	Get(input GetArticleInput) (Article, error)
	IsSlugAvailable(slugArticle string) (bool, error)
	Create(input CreateArticleInput, slugArticle string) (Article, error)
	Update(inputID GetArticleInput, inputData UpdateArticleInput) (Article, error)
	Delete(inputID GetArticleInput, currentUser users.User) (Article, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetAll(userID int) ([]Article, error) {
	if userID != 0 {
		articles, err := s.repository.FindByUserID(userID)
		if err != nil {
			return articles, err
		}

		return articles, nil
	}

	articles, err := s.repository.FindAll()
	if err != nil {
		return articles, err
	}

	return articles, nil
}

func (s *service) Get(input GetArticleInput) (Article, error) {
	article, err := s.repository.FindByID(input.ID)
	if err != nil {
		return article, err
	}
	return article, nil
}
func (s *service) IsSlugAvailable(slugArticle string) (bool, error) {
	article, err := s.repository.FindBySlug(slugArticle)

	if err != nil {
		return false, err
	}

	if article.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) Create(input CreateArticleInput, slugArticle string) (Article, error) {
	article := Article{}
	article.Title = input.Title
	article.Content = input.Content
	datePublished, _ := time.ParseInLocation("2006-01-02 03:04:05", input.PublishedAt, time.Local)
	article.PublishedAt = datePublished
	article.AuthorID = input.User.ID
	article.Slug = slugArticle

	new, err := s.repository.Save(article)
	if err != nil {
		return new, err
	}

	return new, nil
}

func (s *service) Update(inputID GetArticleInput, inputData UpdateArticleInput) (Article, error) {
	article, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return article, err
	}

	if article.AuthorID != inputData.User.ID {
		return article, errors.New("not an owner of the article")
	}

	article.Content = inputData.Content
	datePublished, _ := time.ParseInLocation("2006-01-02 03:04:05", inputData.PublishedAt, time.Local)
	article.PublishedAt = datePublished

	update, err := s.repository.Update(article)
	if err != nil {
		return update, err
	}

	return update, nil
}

func (s *service) Delete(inputID GetArticleInput, currentUser users.User) (Article, error) {
	article, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return article, err
	}

	if article.AuthorID != currentUser.ID {
		return article, errors.New("not an owner of the article")
	}

	delete, err := s.repository.Delete(article)
	if err != nil {
		return delete, err
	}

	return delete, nil
}
