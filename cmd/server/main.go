package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"html/template"
	"news_service/internal/handler"
	"news_service/internal/repository/mongodb"
	"news_service/internal/service"
)

func main() {
	// Initialize MongoDB connection
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

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	// Initialize dependencies
	database := os.Getenv("MONGODB_DATABASE")
	if database == "" {
		database = "news_service"
	}

	newsRepo := mongodb.NewNewsRepository(client, database)
	newsService := service.NewNewsService(newsRepo)
	newsHandler := handler.NewNewsHandler(newsService)

	// Initialize Gin router
	router := gin.Default()

	// Register custom template functions
	funcMap := template.FuncMap{
		"subtract": func(a, b int) int { return a - b },
		"add":      func(a, b int) int { return a + b },
		"multiply": func(a, b int) int { return a * b },
	}
	tmpl := template.Must(template.New("").Funcs(funcMap).ParseGlob("web/templates/**/*"))
	router.SetHTMLTemplate(tmpl)

	router.Static("/static", "./web/static")

	// Register routes
	newsHandler.RegisterRoutes(router)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
