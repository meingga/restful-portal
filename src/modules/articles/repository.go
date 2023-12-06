package articles

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Article, error)
	FindByUserID(userID int) ([]Article, error)
	FindByID(ID int) (Article, error)
	FindBySlug(slug string) (Article, error)
	Save(article Article) (Article, error)
	Update(article Article) (Article, error)
	Delete(article Article) (Article, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Article, error) {
	var articles []Article

	err := r.db.Find(&articles).Error

	if err != nil {
		return articles, err
	}

	return articles, nil
}
func (r *repository) FindByUserID(userID int) ([]Article, error) {
	var articles []Article

	err := r.db.Where("author_id = ?", userID).Find(&articles).Error

	if err != nil {
		return articles, err
	}

	return articles, nil
}

func (r *repository) FindBySlug(slug string) (Article, error) {
	var articles Article

	err := r.db.Where("slug = ?", slug).Find(&articles).Error

	if err != nil {
		return articles, err
	}

	return articles, nil
}

func (r *repository) FindByID(ID int) (Article, error) {
	var article Article

	err := r.db.Where("id = ?", ID).Find(&article).Error

	if err != nil {
		return article, err
	}

	return article, nil
}

func (r *repository) Save(article Article) (Article, error) {
	err := r.db.Create(&article).Error
	if err != nil {
		return article, err
	}
	return article, nil
}

func (r *repository) Update(article Article) (Article, error) {
	err := r.db.Save(&article).Error

	if err != nil {
		return article, err
	}
	return article, nil
}

func (r *repository) Delete(article Article) (Article, error) {
	err := r.db.Delete(&article).Error

	if err != nil {
		return article, err
	}
	return article, nil
}
