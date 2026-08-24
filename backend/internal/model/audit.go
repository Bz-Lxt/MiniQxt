package model

type AnswerTrace struct {
	ID            uint64  `json:"id" gorm:"primaryKey"`
	TenantID      uint64  `json:"tenant_id" gorm:"index"`
	SessionID     uint64  `json:"session_id" gorm:"index:idx_trace_sess_seq"`
	QuestionID    uint64  `json:"question_id"`
	FromOptionID  *uint64 `json:"from_option_id"`
	ToOptionID    *uint64 `json:"to_option_id"`
	AnswerText    string  `json:"answer_text" gorm:"type:text"`
	OccurredAt    Time    `json:"occurred_at"`
	Seq           int     `json:"seq" gorm:"index:idx_trace_sess_seq"`
	CreatedAt     Time    `json:"created_at"`
}

type AntiCheatEvent struct {
	ID          uint64 `json:"id" gorm:"primaryKey"`
	TenantID    uint64 `json:"tenant_id" gorm:"index"`
	SessionID   uint64 `json:"session_id" gorm:"index"`
	EventType   string `json:"event_type" gorm:"size:32"` // blur / hidden / focus
	OccurredAt  Time   `json:"occurred_at"`
	DurationMS  int    `json:"duration_ms"`
	CreatedAt   Time   `json:"created_at"`
}

type AuditFlag struct {
	ID         uint64 `json:"id" gorm:"primaryKey"`
	TenantID   uint64 `json:"tenant_id" gorm:"index"`
	SessionID  uint64 `json:"session_id" gorm:"index"`
	UserID     uint64 `json:"user_id"`
	ExamID     uint64 `json:"exam_id"`
	RuleCode   string `json:"rule_code" gorm:"size:32"`
	Title      string `json:"title" gorm:"size:128"`
	Evidence   string `json:"evidence" gorm:"type:text"`
	Severity   string `json:"severity" gorm:"size:16"`
	CreatedAt  Time   `json:"created_at"`
}

type GraderMetric struct {
	ID           uint64 `json:"id" gorm:"primaryKey"`
	Name         string `json:"name" gorm:"size:64;uniqueIndex"`
	Value        int64  `json:"value"`
	UpdatedAt    Time   `json:"updated_at"`
}
