package model

type Tenant struct {
	ID        uint64 `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"size:128;not null"`
	Code      string `json:"code" gorm:"size:64;uniqueIndex;not null"`
	Status    string `json:"status" gorm:"size:16;default:active"`
	BlurWarn  int    `json:"blur_warn" gorm:"default:3"`
	BlurForce int    `json:"blur_force" gorm:"default:8"`
	PassScore int    `json:"pass_score" gorm:"default:60"`
	CreatedAt Time   `json:"created_at"`
	UpdatedAt Time   `json:"updated_at"`
}

