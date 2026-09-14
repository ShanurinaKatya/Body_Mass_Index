package ds

import "time"

type Like struct {
	ID            uint              `gorm:"primaryKey"`
	UserID        uint              `gorm:"not null;uniqueIndex:idx_likes_user_calculation;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	CalculationID uint              `gorm:"not null;uniqueIndex:idx_likes_user_calculation;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	CreatedAt     time.Time         `gorm:"type:timestamptz;default:now()"`
	User          User              `gorm:"foreignKey:UserID"`
	Category      PatientsCategorie `gorm:"foreignKey:CalculationID"`
}
