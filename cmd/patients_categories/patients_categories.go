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

	rep, err := repository.New(&repository.RepositorySettings{
		PostgresDSN:     dsn.FromEnv(),
		MinioEndpoint:   envOr("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey:  envOr("MINIO_ACCESS_KEY", "root"),
		MinioSecretKey:  envOr("MINIO_SECRET_KEY", "rootpassword"),
		MinioBucketName: envOr("MINIO_BUCKET_NAME", "body-mass-index"),
	})
	if err != nil {
		logrus.Fatal("Repository init error:", err)
	}

	h := handler.NewHandler(rep)

	router := gin.Default()
	app := pkg.NewApp(cfg, router, h)
	app.RunApp()
}

func envOr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}
