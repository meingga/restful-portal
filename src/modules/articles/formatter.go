package articles

import "time"

type ArticleFormatter struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	AuthorID    int       `json:"author_id"`
	Slug        string    `json:"slug"`
	PublishedAt time.Time `json:"published_at"`
}

func FormatArticles(articles []Article) []ArticleFormatter {
	articlesFormatter := []ArticleFormatter{}

	for _, article := range articles {
		articleFormatter := FormatArticle(article)
		articlesFormatter = append(articlesFormatter, articleFormatter)
	}

	return articlesFormatter
}

func FormatArticle(article Article) ArticleFormatter {
	formatter := ArticleFormatter{
		ID:          article.ID,
		Title:       article.Title,
		Content:     article.Content,
		AuthorID:    article.AuthorID,
		PublishedAt: article.PublishedAt,
		Slug:        article.Slug,
	}

	return formatter
}
