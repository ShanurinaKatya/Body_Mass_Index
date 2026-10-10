package ds

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"type:varchar(50);not null;unique" json:"username"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
}

type PatientsCategorie struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"type:varchar(20);not null;default:'draft';check:status_in,status in ('draft','published','deleted')" json:"-"`
	ImageURL    string    `gorm:"type:varchar(255);not null;default:'/static/image/default.jpg'" json:"image_url"`
	VideoURL    string    `gorm:"type:varchar(255);not null;default:'/static/image/Default.mp4'" json:"video_url"`
	Age         int       `json:"age"`
	Weight      float64   `gorm:"type:numeric(5,2)" json:"weight"`
	CreatedAt   time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;default:now()" json:"updated_at"`
	UserID      uint      `gorm:"not null;index;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"-"`
	User        User      `gorm:"foreignKey:UserID" json:"-"`
}

func (PatientsCategorie) TableName() string {
	return "patients_categories"
}

type Like struct {
	ID                   uint              `gorm:"primaryKey" json:"id"`
	UserID               uint              `gorm:"not null;uniqueIndex:idx_likes_user_patients_categories;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	PatientsCategoriesID uint              `gorm:"not null;uniqueIndex:idx_likes_user_patients_categories;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	User                 User              `gorm:"foreignKey:UserID"`
	Category             PatientsCategorie `gorm:"foreignKey:PatientsCategoriesID"`
}

type UserSerializer struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

func SerializeUser(u User) UserSerializer {
	return UserSerializer{
		ID:       u.ID,
		Username: u.Username,
	}
}

type PatientsCategoriesSerializer struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	VideoURL    string    `json:"video_url"`
	Age         int       `json:"age"`
	Weight      float64   `json:"weight"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Likes       int       `json:"likes"`
	CreatorID   uint      `json:"creator_id"`
}

func SerializePatientsCategorie(pc PatientsCategorie, likes int) PatientsCategoriesSerializer {
	return PatientsCategoriesSerializer{
		ID:          pc.ID,
		Name:        pc.Name,
		Description: pc.Description,
		ImageURL:    pc.ImageURL,
		VideoURL:    pc.VideoURL,
		Age:         pc.Age,
		Weight:      pc.Weight,
		CreatedAt:   pc.CreatedAt,
		UpdatedAt:   pc.UpdatedAt,
		Likes:       likes,
		CreatorID:   pc.UserID,
	}
}

// PatientsCategoriesTileSerializer — модель плитки: базовая модель услуги
// плюс признак «карточка создана текущим пользователем».
type PatientsCategoriesTileSerializer struct {
	PatientsCategoriesSerializer
	IsCreator int `json:"is_creator"`
}

func SerializePatientsCategorieTile(pc PatientsCategorie, likes, isCreator int) PatientsCategoriesTileSerializer {
	return PatientsCategoriesTileSerializer{
		PatientsCategoriesSerializer: SerializePatientsCategorie(pc, likes),
		IsCreator:                    isCreator,
	}
}

func SerializePatientsCategoriesTiles(items []PatientsCategorie, likes map[int]int, currentUserID uint) []PatientsCategoriesTileSerializer {
	res := make([]PatientsCategoriesTileSerializer, 0, len(items))
	for _, it := range items {
		isCreator := 0
		if it.UserID == currentUserID {
			isCreator = 1
		}
		res = append(res, SerializePatientsCategorieTile(it, likes[int(it.ID)], isCreator))
	}
	return res
}

// PatientsCategoriesFeedSerializer — модель ленты: базовая модель услуги
// плюс признак «текущий пользователь лайкнул эту карточку».
type PatientsCategoriesFeedSerializer struct {
	PatientsCategoriesSerializer
	IsLiked int `json:"is_liked"`
}

func SerializePatientsCategorieFeed(pc PatientsCategorie, likes, isLiked int) PatientsCategoriesFeedSerializer {
	return PatientsCategoriesFeedSerializer{
		PatientsCategoriesSerializer: SerializePatientsCategorie(pc, likes),
		IsLiked:                      isLiked,
	}
}
