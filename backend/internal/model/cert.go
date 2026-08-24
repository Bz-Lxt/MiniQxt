package model

type CertProgram struct {
	ID            uint64  `json:"id" gorm:"primaryKey"`
	TenantID      uint64  `json:"tenant_id" gorm:"index"`
	Name          string  `json:"name" gorm:"size:128"`
	PositionID    uint64  `json:"position_id"`
	Level         int     `json:"level"`
	RequireCourse uint64  `json:"require_course"`
	RequireExam   uint64  `json:"require_exam"`
	MinScore      float64 `json:"min_score"`
	ValidDays     int     `json:"valid_days" gorm:"default:365"`
	CreatedAt     Time    `json:"created_at"`
}

type Certificate struct {
	ID         uint64 `json:"id" gorm:"primaryKey"`
	TenantID   uint64 `json:"tenant_id" gorm:"index"`
	ProgramID  uint64 `json:"program_id"`
	UserID     uint64 `json:"user_id" gorm:"index"`
	No         string `json:"no" gorm:"size:64;uniqueIndex"`
	Status     string `json:"status" gorm:"size:16"`
	IssuedAt   Time   `json:"issued_at"`
	ExpireAt   Time   `json:"expire_at"`
	Score      float64 `json:"score"`
	Program    *CertProgram `json:"program,omitempty" gorm:"foreignKey:ProgramID"`
}
