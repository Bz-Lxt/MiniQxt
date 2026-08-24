package service

import (
	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/pkg/apperr"
	"github.com/miniqxt/backend/internal/pkg/passwd"
	"github.com/miniqxt/backend/internal/repo"
	"github.com/miniqxt/backend/internal/timeutil"
)

func (a *App) ListDepts(tid uint64) ([]model.Department, error) {
	var list []model.Department
	err := repo.Tenant(a.DB, tid).Order("id asc").Find(&list).Error
	return list, err
}

func (a *App) CreateDept(tid uint64, name string, parent uint64) (model.Department, error) {
	if name == "" {
		return model.Department{}, apperr.Validation.With("部门名称必填")
	}
	d := model.Department{TenantID: tid, Name: name, ParentID: parent, CreatedAt: model.Time(timeutil.Now())}
	if err := a.DB.Create(&d).Error; err != nil {
		return d, apperr.Conflict.With("部门已存在")
	}
	return d, nil
}

func (a *App) PatchDept(tid, id uint64, name string) error {
	res := repo.Tenant(a.DB, tid).Model(&model.Department{}).Where("id = ?", id).Update("name", name)
	if res.RowsAffected == 0 {
		return apperr.NotFound
	}
	return nil
}

func (a *App) DeleteDept(tid, id uint64) error {
	res := repo.Tenant(a.DB, tid).Delete(&model.Department{}, id)
	if res.RowsAffected == 0 {
		return apperr.NotFound
	}
	return nil
}

func (a *App) ListPositions(tid uint64) ([]model.Position, error) {
	var list []model.Position
	err := repo.Tenant(a.DB, tid).Order("level asc").Find(&list).Error
	return list, err
}

func (a *App) CreatePosition(tid uint64, name string, level int) (model.Position, error) {
	if name == "" {
		return model.Position{}, apperr.Validation.With("岗位名称必填")
	}
	if level < 1 {
		level = 1
	}
	p := model.Position{TenantID: tid, Name: name, Level: level, CreatedAt: model.Time(timeutil.Now())}
	if err := a.DB.Create(&p).Error; err != nil {
		return p, apperr.Conflict.With("岗位已存在")
	}
	return p, nil
}

func (a *App) PatchPosition(tid, id uint64, name string, level int) error {
	res := repo.Tenant(a.DB, tid).Model(&model.Position{}).Where("id = ?", id).Updates(map[string]any{"name": name, "level": level})
	if res.RowsAffected == 0 {
		return apperr.NotFound
	}
	return nil
}

func (a *App) DeletePosition(tid, id uint64) error {
	res := repo.Tenant(a.DB, tid).Delete(&model.Position{}, id)
	if res.RowsAffected == 0 {
		return apperr.NotFound
	}
	return nil
}

type EmployeeIn struct {
	Username    string  `json:"username"`
	Password    string  `json:"password"`
	DisplayName string  `json:"display_name"`
	Role        string  `json:"role"`
	DeptID      *uint64 `json:"dept_id"`
	PositionID  *uint64 `json:"position_id"`
}

func (a *App) ListEmployees(tid uint64, page, size int, q string) ([]model.User, int64, error) {
	tx := repo.Tenant(a.DB.Model(&model.User{}), tid).Where("role <> ?", model.RolePlatform)
	if q != "" {
		like := "%" + q + "%"
		tx = tx.Where("username LIKE ? OR display_name LIKE ?", like, like)
	}
	var total int64
	tx.Count(&total)
	var list []model.User
	err := repo.Page(tx.Preload("Department").Preload("Position"), page, size).Order("id desc").Find(&list).Error
	return list, total, err
}

func (a *App) CreateEmployee(tid uint64, in EmployeeIn) (model.User, error) {
	if in.Username == "" || in.Password == "" || in.DisplayName == "" {
		return model.User{}, apperr.Validation.With("用户名、密码、姓名必填")
	}
	if in.Role == "" {
		in.Role = model.RoleEmployee
	}
	if in.Role == model.RolePlatform {
		return model.User{}, apperr.Validation.With("不可创建平台账号")
	}
	hash, err := passwd.Hash(in.Password)
	if err != nil {
		return model.User{}, apperr.Internal
	}
	u := model.User{
		TenantID: tid, Username: in.Username, PasswordHash: hash,
		DisplayName: in.DisplayName, Role: in.Role, DeptID: in.DeptID, PositionID: in.PositionID,
		Status: model.StatusActive, CreatedAt: model.Time(timeutil.Now()),
	}
	if err := a.DB.Create(&u).Error; err != nil {
		return u, apperr.Conflict.With("用户名已存在")
	}
	return u, nil
}

func (a *App) PatchEmployee(tid, id uint64, in EmployeeIn) error {
	updates := map[string]any{}
	if in.DisplayName != "" {
		updates["display_name"] = in.DisplayName
	}
	if in.Role != "" && in.Role != model.RolePlatform {
		updates["role"] = in.Role
	}
	if in.DeptID != nil {
		updates["dept_id"] = *in.DeptID
	}
	if in.PositionID != nil {
		updates["position_id"] = *in.PositionID
	}
	if in.Password != "" {
		hash, err := passwd.Hash(in.Password)
		if err != nil {
			return apperr.Internal
		}
		updates["password_hash"] = hash
	}
	res := repo.Tenant(a.DB, tid).Model(&model.User{}).Where("id = ?", id).Updates(updates)
	if res.RowsAffected == 0 {
		return apperr.NotFound
	}
	return nil
}

func (a *App) DeleteEmployee(tid, id uint64) error {
	res := repo.Tenant(a.DB, tid).Model(&model.User{}).Where("id = ?", id).Update("status", model.StatusDisabled)
	if res.RowsAffected == 0 {
		return apperr.NotFound
	}
	return nil
}
