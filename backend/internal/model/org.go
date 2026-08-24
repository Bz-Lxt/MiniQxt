package model

type Department struct {
	ID        uint64 `json:"id" gorm:"primaryKey"`
	TenantID  uint64 `json:"tenant_id" gorm:"index;uniqueIndex:uk_dept_name"`
	Name      string `json:"name" gorm:"size:128;uniqueIndex:uk_dept_name"`
	ParentID  uint64 `json:"parent_id"`
	CreatedAt Time   `json:"created_at"`
}

type Position struct {
	ID        uint64 `json:"id" gorm:"primaryKey"`
	TenantID  uint64 `json:"tenant_id" gorm:"index;uniqueIndex:uk_pos_name"`
	Name      string `json:"name" gorm:"size:128;uniqueIndex:uk_pos_name"`
	Level     int    `json:"level" gorm:"default:1"`
	CreatedAt Time   `json:"created_at"`
}
