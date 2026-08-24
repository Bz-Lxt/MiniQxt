package model

type Course struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	TenantID    uint64    `json:"tenant_id" gorm:"index"`
	Title       string    `json:"title" gorm:"size:256;not null"`
	Summary     string    `json:"summary" gorm:"type:text"`
	Cover       string    `json:"cover" gorm:"size:256"`
	Required    bool      `json:"required"`
	CreatedAt   Time      `json:"created_at"`
	Chapters    []Chapter `json:"chapters,omitempty" gorm:"foreignKey:CourseID"`
}

type Chapter struct {
	ID          uint64  `json:"id" gorm:"primaryKey"`
	TenantID    uint64  `json:"tenant_id" gorm:"index"`
	CourseID    uint64  `json:"course_id" gorm:"index"`
	Title       string  `json:"title" gorm:"size:256"`
	VideoFile   string  `json:"video_file" gorm:"size:256"`
	DurationSec int     `json:"duration_sec"`
	SortNo      int     `json:"sort_no"`
	CreatedAt   Time    `json:"created_at"`
}

type CourseAssignment struct {
	ID         uint64 `json:"id" gorm:"primaryKey"`
	TenantID   uint64 `json:"tenant_id" gorm:"index"`
	CourseID   uint64 `json:"course_id" gorm:"index"`
	Scope      string `json:"scope" gorm:"size:16"` // all / dept / user
	TargetID   uint64 `json:"target_id"`
	CreatedAt  Time   `json:"created_at"`
}

type LearningProgress struct {
	ID            uint64  `json:"id" gorm:"primaryKey"`
	TenantID      uint64  `json:"tenant_id" gorm:"uniqueIndex:uk_progress"`
	UserID        uint64  `json:"user_id" gorm:"uniqueIndex:uk_progress"`
	ChapterID     uint64  `json:"chapter_id" gorm:"uniqueIndex:uk_progress"`
	CourseID      uint64  `json:"course_id" gorm:"index"`
	WatchedSec    int     `json:"watched_sec"`
	PositionSec   int     `json:"position_sec"`
	Percent       float64 `json:"percent"`
	LastReportAt  Time    `json:"last_report_at"`
	Completed     bool    `json:"completed"`
}
