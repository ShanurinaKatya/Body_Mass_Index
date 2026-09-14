package main

import (
	"os"

	"patients_categories/internal/app/config"
	"patients_categories/internal/app/dsn"
	"patients_categories/internal/app/handler"
	"patients_categories/internal/app/repository"
	"patients_categories/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Info(".env file not found, using defaults")
	}

	cfg, err := config.NewConfig()
	if err != nil {
		logrus.Fatal("Config init error:", err)
	}

	rep, err := repository.New(dsn.FromEnv())
	if err != nil {
		logrus.Fatal("Repository init error:", err)
	}

	minioURL := "http://localhost:9000/body-mass-index"
	if v := os.Getenv("MINIO_URL"); v != "" {
		minioURL = v
	}

	h := handler.NewHandler(rep, minioURL)

	router := gin.Default()
	app := pkg.NewApp(cfg, router, h)
	app.RunApp()
}
