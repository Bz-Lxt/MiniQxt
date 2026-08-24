package model

type Exam struct {
	ID          uint64 `json:"id" gorm:"primaryKey"`
	TenantID    uint64 `json:"tenant_id" gorm:"index"`
	PaperID     uint64 `json:"paper_id"`
	Title       string `json:"title" gorm:"size:256"`
	StartAt     Time   `json:"start_at"`
	EndAt       Time   `json:"end_at"`
	DurationSec int    `json:"duration_sec"`
	PassScore   float64 `json:"pass_score"`
	MaxAttempts int    `json:"max_attempts" gorm:"default:1"`
	Status      string `json:"status" gorm:"size:16;default:published"`
	CreatedAt   Time   `json:"created_at"`
	Paper       *Paper `json:"paper,omitempty" gorm:"foreignKey:PaperID"`
}

type ExamAssignment struct {
	ID        uint64 `json:"id" gorm:"primaryKey"`
	TenantID  uint64 `json:"tenant_id" gorm:"index"`
	ExamID    uint64 `json:"exam_id" gorm:"index"`
	Scope     string `json:"scope" gorm:"size:16"`
	TargetID  uint64 `json:"target_id"`
	CreatedAt Time   `json:"created_at"`
}

type ExamSession struct {
	ID            uint64  `json:"id" gorm:"primaryKey"`
	TenantID      uint64  `json:"tenant_id" gorm:"uniqueIndex:uk_session_attempt"`
	ExamID        uint64  `json:"exam_id" gorm:"uniqueIndex:uk_session_attempt"`
	UserID        uint64  `json:"user_id" gorm:"uniqueIndex:uk_session_attempt"`
	Attempt       int     `json:"attempt" gorm:"uniqueIndex:uk_session_attempt"`
	ShuffleSeed   int64   `json:"shuffle_seed"`
	StartedAt     Time    `json:"started_at"`
	EndsAt        Time    `json:"ends_at"`
	Status        string  `json:"status" gorm:"size:24;index"`
	BlurCount     int     `json:"blur_count"`
	HeartbeatN    int     `json:"heartbeat_n"`
	Suspicious    bool    `json:"suspicious"`
	Integrity     string  `json:"integrity" gorm:"size:32;default:ok"`
	ForceReason   string  `json:"force_reason" gorm:"size:128"`
	CreatedAt     Time    `json:"created_at"`
}
