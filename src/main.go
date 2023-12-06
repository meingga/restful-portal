package main

import (
	"log"
	helperJWT "restful-portal/src/helpers/jwt"
	"restful-portal/src/helpers/mysql"
	"restful-portal/src/middleware"
	"restful-portal/src/modules/articles"
	"restful-portal/src/modules/users"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := mysql.IsConnected()
	if err != nil {
		log.Fatal(err.Error())
	}

	authService := helperJWT.NewService()

	usersRepository := users.NewRepository(db)
	usersService := users.NewService(usersRepository)
	usersHandler := users.NewUserHandler(usersService, authService)

	articlesRepository := articles.NewRepository(db)
	articlesService := articles.NewService(articlesRepository)
	articlesHandler := articles.NewHandler(articlesService)

	router := gin.Default()

	router.Use(cors.Default())
	api := router.Group("api/v1")

	api.POST("/users", usersHandler.Register)
	api.POST("/login", usersHandler.Login)
	api.GET("/refresh-token", usersHandler.RefreshToken)
	api.GET("/users/fetch", middleware.AuthMiddleware(authService, usersService), usersHandler.FetchUser)
	api.PUT("/users/:id", middleware.AuthMiddleware(authService, usersService), usersHandler.Update)

	api.GET("/articles", middleware.AuthMiddleware(authService, usersService), articlesHandler.GetArticles)
	api.GET("/articles/:id", middleware.AuthMiddleware(authService, usersService), articlesHandler.GetArticle)
	api.POST("/articles", middleware.AuthMiddleware(authService, usersService), articlesHandler.CreateArticle)
	api.PUT("/articles/:id", middleware.AuthMiddleware(authService, usersService), articlesHandler.UpdateArticle)
	api.DELETE("/articles/:id", middleware.AuthMiddleware(authService, usersService), articlesHandler.DeleteArticle)

	router.Run(":8383")
}
