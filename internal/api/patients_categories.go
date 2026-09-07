package api

import (
	"log"
	"patients_categories/internal/app/handler"
	"patients_categories/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		log.Fatal("Repository init error:", err)
	}

	minioURL := "http://localhost:9000/body-mass-index"
	h := handler.NewHandler(repo, minioURL)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/", h.FeedHandler)
	r.GET("/feed", h.FeedHandler)
	r.GET("/add", h.AddHandler)
	r.GET("/cards", h.CardsHandler)

	log.Println("Server running on http://localhost:8080")
	r.Run(":8080")

	log.Println("Server down")
}
