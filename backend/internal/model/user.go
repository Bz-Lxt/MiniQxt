package model

type User struct {
	ID           uint64  `json:"id" gorm:"primaryKey"`
	TenantID     uint64  `json:"tenant_id" gorm:"index"`
	Username     string  `json:"username" gorm:"size:128;uniqueIndex;not null"`
	PasswordHash string  `json:"-" gorm:"size:128;not null"`
	DisplayName  string  `json:"display_name" gorm:"size:128"`
	Role         string  `json:"role" gorm:"size:32;index"`
	DeptID       *uint64 `json:"dept_id"`
	PositionID   *uint64 `json:"position_id"`
	Status       string  `json:"status" gorm:"size:16;default:active"`
	CreatedAt    Time    `json:"created_at"`
	UpdatedAt    Time    `json:"updated_at"`

	Department *Department `json:"department,omitempty" gorm:"foreignKey:DeptID"`
	Position   *Position   `json:"position,omitempty" gorm:"foreignKey:PositionID"`
}

type PublicUser struct {
	ID          uint64 `json:"id"`
	TenantID    uint64 `json:"tenant_id,omitempty"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
	DeptID      *uint64 `json:"dept_id"`
	PositionID  *uint64 `json:"position_id"`
}

func (u User) Public() PublicUser {
	return PublicUser{
		ID: u.ID, TenantID: u.TenantID, Username: u.Username,
		DisplayName: u.DisplayName, Role: u.Role, DeptID: u.DeptID, PositionID: u.PositionID,
	}
}
