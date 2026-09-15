package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"patients_categories/internal/app/ds"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DraftStatus     = "draft"
	PublishedStatus = "published"
	DeletedStatus   = "deleted"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

func (r *Repository) GetPublished() ([]ds.PatientsCategorie, error) {
	var calcs []ds.PatientsCategorie
	err := r.db.Where("status = ?", PublishedStatus).Order("id").Find(&calcs).Error
	return calcs, err
}

func (r *Repository) GetNext(id uint) (ds.PatientsCategorie, error) {
	published, err := r.GetPublished()
	if err != nil {
		return ds.PatientsCategorie{}, err
	}
	if len(published) == 0 {
		return ds.PatientsCategorie{}, fmt.Errorf("published categories not found")
	}
	for i, c := range published {
		if c.ID == id {
			return published[(i+1)%len(published)], nil
		}
	}
	return ds.PatientsCategorie{}, fmt.Errorf("category %d not found", id)
}

func (r *Repository) GetPublishedByID(id uint) (ds.PatientsCategorie, error) {
	var calc ds.PatientsCategorie
	err := r.db.Where("id = ? AND status = ?", id, PublishedStatus).First(&calc).Error
	return calc, err
}

func (r *Repository) GetDraftByUser(userID uint) (*ds.PatientsCategorie, error) {
	var calc ds.PatientsCategorie
	err := r.db.Where("user_id = ? AND status = ?", userID, DraftStatus).First(&calc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &calc, nil
}

func (r *Repository) GetPublishedByName(name string) (*ds.PatientsCategorie, error) {
	var calc ds.PatientsCategorie
	err := r.db.Where("status = ? AND name = ?", PublishedStatus, name).First(&calc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &calc, nil
}

func (r *Repository) GetPublishedByAge(age int) ([]ds.PatientsCategorie, error) {
	var calcs []ds.PatientsCategorie
	err := r.db.Where("status = ? AND age = ?", PublishedStatus, age).Order("id").Find(&calcs).Error
	return calcs, err
}

func (r *Repository) CreateCategory(name string, userID uint) (ds.PatientsCategorie, error) {
	calc := ds.PatientsCategorie{
		Name:   name,
		Status: DraftStatus,
		UserID: userID,
	}
	err := r.db.Select("Name", "Status", "UserID").Create(&calc).Error
	return calc, err
}

func (r *Repository) PrefillDraft(id uint, src ds.PatientsCategorie) error {
	return r.db.Model(&ds.PatientsCategorie{}).
		Where("id = ? AND status = ?", id, DraftStatus).
		Updates(map[string]interface{}{
			"description": src.Description,
			"image_url":   src.ImageURL,
			"video_url":   src.VideoURL,
			"age":         src.Age,
			"gender":      src.Gender,
			"weight":      src.Weight,
			"height":      src.Height,
			"updated_at":  time.Now(),
		}).Error
}

func (r *Repository) PublishCategory(id uint, name, description string, age int, gender string, weight float64, height int) error {
	return r.db.Model(&ds.PatientsCategorie{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":        name,
			"description": description,
			"age":         age,
			"gender":      gender,
			"weight":      weight,
			"height":      height,
			"status":      PublishedStatus,
			"updated_at":  time.Now(),
		}).Error
}

func (r *Repository) DeleteCategory(id uint) error {
	row := r.db.Raw("SELECT id, name FROM patients_categories WHERE id = $1", id).Row()
	var calcID uint
	var name string
	if err := row.Scan(&calcID, &name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("category %d not found", id)
		}
		return fmt.Errorf("error selecting category: %w", err)
	}

	err := r.db.Exec("UPDATE patients_categories SET status = ?, updated_at = NOW() WHERE id = ?",
		DeletedStatus, calcID).Error
	if err != nil {
		return fmt.Errorf("error deleting category: %w", err)
	}
	return nil
}

func (r *Repository) GetLikesCounts() (map[int]int, error) {
	var rows []struct {
		CalculationID uint
		Cnt           int
	}
	err := r.db.Model(&ds.Like{}).
		Select("calculation_id, COUNT(*) AS cnt").
		Group("calculation_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	res := make(map[int]int, len(rows))
	for _, row := range rows {
		res[int(row.CalculationID)] = row.Cnt
	}
	return res, nil
}

func (r *Repository) GetLikesCount(id uint) (int, error) {
	var count int64
	err := r.db.Model(&ds.Like{}).Where("calculation_id = ?", id).Count(&count).Error
	return int(count), err
}