package main

import (
	"patients_categories/internal/app/ds"
	"patients_categories/internal/app/dsn"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Info(".env file not found, using defaults")
	}

	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		logrus.Fatal("failed to connect database:", err)
	}

	if err := db.AutoMigrate(&ds.User{}, &ds.PatientsCategorie{}, &ds.Like{}); err != nil {
		logrus.Fatal("migration failed:", err)
	}

	logrus.Info("migration completed")
}
