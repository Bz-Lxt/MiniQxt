package model

type Submission struct {
	ID              uint64  `json:"id" gorm:"primaryKey"`
	TenantID        uint64  `json:"tenant_id" gorm:"index"`
	SessionID       uint64  `json:"session_id" gorm:"uniqueIndex"`
	ExamID          uint64  `json:"exam_id" gorm:"index"`
	UserID          uint64  `json:"user_id" gorm:"index"`
	Status          string  `json:"status" gorm:"size:24;index"`
	ObjectiveScore  float64 `json:"objective_score"`
	SubjectiveScore float64 `json:"subjective_score"`
	TotalScore      float64 `json:"total_score"`
	Pass            bool    `json:"pass"`
	Integrity       string  `json:"integrity" gorm:"size:32"`
	ErrorMsg        string  `json:"error_msg,omitempty" gorm:"size:512"`
	QueuedAt        Time    `json:"queued_at"`
	GradedAt        *Time   `json:"graded_at"`
	CreatedAt       Time    `json:"created_at"`
	Answers         []SubmissionAnswer `json:"answers,omitempty" gorm:"foreignKey:SubmissionID"`
}

type SubmissionAnswer struct {
	ID           uint64  `json:"id" gorm:"primaryKey"`
	SubmissionID uint64  `json:"submission_id" gorm:"index"`
	QuestionID   uint64  `json:"question_id"`
	OptionIDs    string  `json:"option_ids" gorm:"size:256"` // comma-separated stable ids
	AnswerText   string  `json:"answer_text" gorm:"type:text"`
	AutoScore    float64 `json:"auto_score"`
	ManualScore  *float64 `json:"manual_score"`
	Comment      string  `json:"comment" gorm:"size:512"`
	IsEssay      bool    `json:"is_essay"`
}
