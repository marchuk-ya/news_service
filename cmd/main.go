package main

import (
	"context"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"news_service/internal/handler"
	"news_service/internal/repository/mongodb"
	"news_service/internal/service"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	database := os.Getenv("MONGODB_DATABASE")
	if database == "" {
		database = "news_service"
	}

	newsRepo := mongodb.NewNewsRepository(client, database)
	newsService := service.NewNewsService(newsRepo)
	newsHandler := handler.NewNewsHandler(newsService)

	router := gin.Default()

	funcMap := template.FuncMap{
		"subtract": func(a, b int) int { return a - b },
		"add":      func(a, b int) int { return a + b },
		"multiply": func(a, b int) int { return a * b },
	}
	router.SetFuncMap(funcMap)
	router.LoadHTMLGlob("web/templates/**/*")

	router.Static("/static", "./web/static")

	newsHandler.RegisterRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
