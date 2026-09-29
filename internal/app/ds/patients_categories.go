package ds

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"type:varchar(50);not null;unique" json:"username"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	IsModerator  bool   `gorm:"type:boolean;default:false" json:"is_moderator"`
}

type PatientsCategorie struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Status      string `gorm:"type:varchar(20);not null;default:'draft';check:status_in,status in ('draft','published','deleted')" json:"-"`
	ImageURL    string `gorm:"type:varchar(255);not null;default:'/static/image/default.jpg'" json:"image_url"`
	VideoURL    string `gorm:"type:varchar(255);not null;default:'/static/image/Default.mp4'" json:"video_url"`
	Age         int    `json:"age"`
	Gender      string `gorm:"type:varchar(10);check:gender_in,gender in ('Male','Female')" json:"gender"`
	Weight      float64 `gorm:"type:numeric(5,2)" json:"weight"`
	Height      int    `json:"height"`
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
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	IsModerator bool   `json:"is_moderator"`
}

func SerializeUser(u User) UserSerializer {
	return UserSerializer{
		ID:          u.ID,
		Username:    u.Username,
		IsModerator: u.IsModerator,
	}
}

type PatientsCategoriesSerializer struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	ImageURL    string         `json:"image_url"`
	VideoURL    string         `json:"video_url"`
	Age         int            `json:"age"`
	Gender      string         `json:"gender"`
	Weight      float64        `json:"weight"`
	Height      int            `json:"height"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Likes       int            `json:"likes"`
	Creator     UserSerializer `json:"creator"`
}

func SerializePatientsCategorie(pc PatientsCategorie, likes int) PatientsCategoriesSerializer {
	return PatientsCategoriesSerializer{
		ID:          pc.ID,
		Name:        pc.Name,
		Description: pc.Description,
		ImageURL:    pc.ImageURL,
		VideoURL:    pc.VideoURL,
		Age:         pc.Age,
		Gender:      pc.Gender,
		Weight:      pc.Weight,
		Height:      pc.Height,
		CreatedAt:   pc.CreatedAt,
		UpdatedAt:   pc.UpdatedAt,
		Likes:       likes,
		Creator:     SerializeUser(pc.User),
	}
}

func SerializePatientsCategories(items []PatientsCategorie, likes map[int]int) []PatientsCategoriesSerializer {
	res := make([]PatientsCategoriesSerializer, 0, len(items))
	for _, it := range items {
		res = append(res, SerializePatientsCategorie(it, likes[int(it.ID)]))
	}
	return res
}