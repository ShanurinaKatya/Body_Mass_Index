package main

import (
	"crypto/sha1"
	"encoding/hex"

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

	// AutoMigrate не удаляет колонки, поэтому устаревшие поля убираем явно
	dropped := []string{
		"ALTER TABLE patients_categories DROP COLUMN IF EXISTS gender",
		"ALTER TABLE patients_categories DROP COLUMN IF EXISTS height",
		"ALTER TABLE users DROP COLUMN IF EXISTS is_moderator",
	}
	for _, stmt := range dropped {
		if err := db.Exec(stmt).Error; err != nil {
			logrus.Fatal("migration failed:", err)
		}
	}

	// Пароли лабы 3 хранились открытым текстом — переводим их на sha1-хэши.
	// Идемпотентно: строка из 40 hex-символов уже является хэшем.
	if err := hashPlaintextPasswords(db); err != nil {
		logrus.Fatal("migration failed:", err)
	}

	logrus.Info("migration completed")
}

func hashPlaintextPasswords(db *gorm.DB) error {
	var users []ds.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}
	for _, u := range users {
		if isSHA1Hex(u.PasswordHash) {
			continue
		}
		if err := db.Model(&ds.User{}).Where("id = ?", u.ID).
			Update("password_hash", sha1Hex(u.PasswordHash)).Error; err != nil {
			return err
		}
	}
	return nil
}

func sha1Hex(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func isSHA1Hex(s string) bool {
	if len(s) != 40 {
		return false
	}
	for _, c := range s {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F') {
			return false
		}
	}
	return true
}
