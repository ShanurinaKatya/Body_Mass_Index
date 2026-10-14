package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"patients_categories/internal/app/ds"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DraftStatus     = "draft"
	PublishedStatus = "published"
	DeletedStatus   = "deleted"
)

var (
	ErrNotFound  = errors.New("patients categorie not found")
	ErrForbidden = errors.New("access to patients categorie is forbidden")
)

type Repository struct {
	db          *gorm.DB
	minio       *minio.Client
	minioURL    string
	minioBucket string
}

type RepositorySettings struct {
	PostgresDSN     string
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucketName string
}

func New(settings *RepositorySettings) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(settings.PostgresDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.New(settings.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(settings.MinioAccessKey, settings.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating minio client: %w", err)
	}

	return &Repository{
		db:          db,
		minio:       minioClient,
		minioURL:    "http://" + settings.MinioEndpoint + "/" + settings.MinioBucketName,
		minioBucket: settings.MinioBucketName,
	}, nil
}

func (r *Repository) GetPublished() ([]ds.PatientsCategorie, error) {
	var calcs []ds.PatientsCategorie
	err := r.db.Preload("User").Where("status = ?", PublishedStatus).Order("id").Find(&calcs).Error
	return calcs, err
}

func (r *Repository) GetPublishedFiltered(age int) ([]ds.PatientsCategorie, error) {
	var calcs []ds.PatientsCategorie
	err := r.db.Preload("User").
		Where("status = ? AND age = ?", PublishedStatus, age).
		Order("id").Find(&calcs).Error
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
	err := r.db.Preload("User").Where("id = ? AND status = ?", id, PublishedStatus).First(&calc).Error
	return calc, err
}

func (r *Repository) GetFirstPublished() (ds.PatientsCategorie, error) {
	var calc ds.PatientsCategorie
	err := r.db.Preload("User").Where("status = ?", PublishedStatus).Order("id").First(&calc).Error
	return calc, err
}

func (r *Repository) GetDraftByUser(userID uint) (*ds.PatientsCategorie, error) {
	var calc ds.PatientsCategorie
	err := r.db.Preload("User").Where("user_id = ? AND status = ?", userID, DraftStatus).First(&calc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &calc, nil
}

func (r *Repository) CreateCategoryDraft(c ds.PatientsCategorie) (ds.PatientsCategorie, error) {
	c.Status = DraftStatus
	err := r.db.Select("Name", "Description", "Status", "Age", "Weight", "UserID").Create(&c).Error
	return c, err
}

func (r *Repository) UpdateAndPublishDraft(id, userID uint, description string, age int, weight float64) error {
	var calc ds.PatientsCategorie
	err := r.db.Where("id = ?", id).First(&calc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: category %d", ErrNotFound, id)
		}
		return err
	}
	if calc.UserID != userID {
		return fmt.Errorf("%w: category %d", ErrForbidden, id)
	}
	if calc.Status != DraftStatus {
		return fmt.Errorf("category %d is not a draft, cannot publish", id)
	}
	return r.db.Model(&ds.PatientsCategorie{}).
		Where("id = ? AND status = ?", id, DraftStatus).
		Updates(map[string]interface{}{
			"description": description,
			"age":         age,
			"weight":      weight,
			"status":      PublishedStatus,
			"updated_at":  time.Now(),
		}).Error
}

func (r *Repository) DeleteCategory(id, userID uint) error {
	row := r.db.Raw("SELECT id, name, user_id FROM patients_categories WHERE id = $1", id).Row()
	var calcID uint
	var name string
	var ownerID uint
	if err := row.Scan(&calcID, &name, &ownerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: category %d", ErrNotFound, id)
		}
		return fmt.Errorf("error selecting category: %w", err)
	}

	if ownerID != userID {
		return fmt.Errorf("%w: category %d", ErrForbidden, id)
	}

	err := r.db.Exec("UPDATE patients_categories SET status = ?, updated_at = NOW() WHERE id = ? AND user_id = ?",
		DeletedStatus, calcID, userID).Error
	if err != nil {
		return fmt.Errorf("error deleting category: %w", err)
	}
	return nil
}

func (r *Repository) UploadCategoryImage(id uint, header *multipart.FileHeader) error {
	filename := "pc_" + strconv.FormatUint(uint64(id), 10) + "_img" + fileExt(header.Filename)
	return r.uploadToCategory(id, filename, header, "image_url")
}

func (r *Repository) UploadCategoryVideo(id uint, header *multipart.FileHeader) error {
	filename := "pc_" + strconv.FormatUint(uint64(id), 10) + "_video" + fileExt(header.Filename)
	return r.uploadToCategory(id, filename, header, "video_url")
}

func fileExt(name string) string {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			return name[i:]
		}
	}
	return ""
}

func (r *Repository) uploadToCategory(id uint, filename string, header *multipart.FileHeader, column string) error {
	file, err := header.Open()
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}
	contentType := http.DetectContentType(buffer)
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("ошибка перемещения по файловому потоку: %w", err)
	}

	if _, err := r.minio.PutObject(context.Background(), r.minioBucket, filename, file, header.Size, minio.PutObjectOptions{
		ContentType: contentType,
	}); err != nil {
		return fmt.Errorf("не удалось добавить объект в хранилище minio: %w", err)
	}

	url := r.minioURL + "/" + filename
	if err := r.db.Model(&ds.PatientsCategorie{}).Where("id = ?", id).Update(column, url).Error; err != nil {
		r.minio.RemoveObject(context.Background(), r.minioBucket, filename, minio.RemoveObjectOptions{})
		return fmt.Errorf("ошибка сохранения пути к файлу: %w", err)
	}
	return nil
}

// IsLiked — лайкнул ли пользователь userID услугу categoryID.
func (r *Repository) IsLiked(categoryID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&ds.Like{}).
		Where("patients_categories_id = ? AND user_id = ?", categoryID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) GetLikesCounts() (map[int]int, error) {
	var rows []struct {
		PatientsCategoriesID uint
		Cnt                  int
	}
	err := r.db.Model(&ds.Like{}).
		Select("patients_categories_id, COUNT(*) AS cnt").
		Group("patients_categories_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	res := make(map[int]int, len(rows))
	for _, row := range rows {
		res[int(row.PatientsCategoriesID)] = row.Cnt
	}
	return res, nil
}

func (r *Repository) GetLikesCount(id uint) (int, error) {
	var count int64
	err := r.db.Model(&ds.Like{}).Where("patients_categories_id = ?", id).Count(&count).Error
	return int(count), err
}

func (r *Repository) SetLike(categoryID, userID uint, liked bool) error {
	if liked {
		var exists int64
		err := r.db.Model(&ds.Like{}).
			Where("patients_categories_id = ? AND user_id = ?", categoryID, userID).
			Count(&exists).Error
		if err != nil {
			return err
		}
		if exists > 0 {
			return nil
		}
		return r.db.Create(&ds.Like{
			UserID:               userID,
			PatientsCategoriesID: categoryID,
		}).Error
	}
	return r.db.Where("patients_categories_id = ? AND user_id = ?", categoryID, userID).
		Delete(&ds.Like{}).Error
}

func (r *Repository) CreateUser(username, password string) (ds.User, error) {
	u := ds.User{
		Username:     username,
		PasswordHash: password,
	}
	err := r.db.Select("Username", "PasswordHash").Create(&u).Error
	return u, err
}

func (r *Repository) GetUserByUsername(username string) (*ds.User, error) {
	var u ds.User
	err := r.db.Where("username = ?", username).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
