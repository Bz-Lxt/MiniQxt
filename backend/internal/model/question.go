package model

type Question struct {
	ID         uint64           `json:"id" gorm:"primaryKey"`
	TenantID   uint64           `json:"tenant_id" gorm:"index"`
	Type       string           `json:"type" gorm:"size:16;index"`
	Stem       string           `json:"stem" gorm:"type:text;not null"`
	Difficulty int              `json:"difficulty" gorm:"default:1"`
	Tags       string           `json:"tags" gorm:"size:256"`
	Score      float64          `json:"score" gorm:"default:5"`
	CreatedAt  Time             `json:"created_at"`
	UpdatedAt  Time             `json:"updated_at"`
	Options    []QuestionOption `json:"options,omitempty" gorm:"foreignKey:QuestionID"`
}

type QuestionOption struct {
	ID         uint64 `json:"id" gorm:"primaryKey"`
	QuestionID uint64 `json:"question_id" gorm:"index"`
	Label      string `json:"label" gorm:"type:text"`
	IsCorrect  bool   `json:"is_correct,omitempty"`
	SortNo     int    `json:"sort_no"`
}

type Paper struct {
	ID                uint64      `json:"id" gorm:"primaryKey"`
	TenantID          uint64      `json:"tenant_id" gorm:"index"`
	Title             string      `json:"title" gorm:"size:256"`
	ShuffleQuestions  bool        `json:"shuffle_questions"`
	ShuffleOptions    bool        `json:"shuffle_options"`
	TotalScore        float64     `json:"total_score"`
	CreatedAt         Time        `json:"created_at"`
	UpdatedAt         Time        `json:"updated_at"`
	Items             []PaperItem `json:"items,omitempty" gorm:"foreignKey:PaperID"`
}

type PaperItem struct {
	ID         uint64   `json:"id" gorm:"primaryKey"`
	PaperID    uint64   `json:"paper_id" gorm:"index"`
	QuestionID uint64   `json:"question_id"`
	Score      float64  `json:"score"`
	GroupName  string   `json:"group_name" gorm:"size:64"`
	SortNo     int      `json:"sort_no"`
	Question   Question `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
}
