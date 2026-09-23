package ds

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"type:varchar(50);not null;unique"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
}

type PatientsCategorie struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	Status      string `gorm:"type:varchar(20);not null;default:'draft';check:status_in,status in ('draft','published','deleted')"`
	ImageURL    string `gorm:"type:varchar(255);not null;default:'/static/image/default.jpg'"`
	VideoURL    string `gorm:"type:varchar(255);not null;default:'/static/image/Default.mp4'"`
	Age         int
	Gender      string  `gorm:"type:varchar(10);check:gender_in,gender in ('Male','Female')"`
	Weight      float64 `gorm:"type:numeric(5,2)"`
	Height      int
	CreatedAt   time.Time `gorm:"type:timestamptz;default:now()"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;default:now()"`
	UserID      uint      `gorm:"not null;index;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	User        User      `gorm:"foreignKey:UserID"`
}

func (PatientsCategorie) TableName() string {
	return "patients_categories"
}

type Like struct {
	ID                  uint              `gorm:"primaryKey"`
	UserID              uint              `gorm:"not null;uniqueIndex:idx_likes_user_patients_categories;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	PatientsCategoriesID uint             `gorm:"not null;uniqueIndex:idx_likes_user_patients_categories;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	User                User              `gorm:"foreignKey:UserID"`
	Category            PatientsCategorie `gorm:"foreignKey:PatientsCategoriesID"`
}