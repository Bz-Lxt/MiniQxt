package service

import (
	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/pkg/apperr"
	"github.com/miniqxt/backend/internal/timeutil"
)

type TenantPatch struct {
	Name      *string
	Status    *string
	BlurWarn  *int
	BlurForce *int
	PassScore *int
}

func (a *App) ListTenants() ([]model.Tenant, error) {
	var list []model.Tenant
	err := a.DB.Order("id asc").Find(&list).Error
	return list, err
}

func (a *App) CreateTenant(name, code string) (model.Tenant, error) {
	if name == "" || code == "" {
		return model.Tenant{}, apperr.Validation.With("名称与编码必填")
	}
	t := model.Tenant{
		Name: name, Code: code, Status: model.StatusActive,
		BlurWarn: 3, BlurForce: 8, PassScore: 60, CreatedAt: model.Time(timeutil.Now()),
	}
	if err := a.DB.Create(&t).Error; err != nil {
		return t, apperr.Conflict.With("租户编码已存在")
	}
	return t, nil
}

func (a *App) PatchTenant(id uint64, p TenantPatch) (model.Tenant, error) {
	var t model.Tenant
	if err := a.DB.First(&t, id).Error; err != nil {
		return t, apperr.NotFound
	}
	if p.Name != nil {
		t.Name = *p.Name
	}
	if p.Status != nil {
		t.Status = *p.Status
	}
	if p.BlurWarn != nil {
		t.BlurWarn = *p.BlurWarn
	}
	if p.BlurForce != nil {
		t.BlurForce = *p.BlurForce
	}
	if p.PassScore != nil {
		t.PassScore = *p.PassScore
	}
	t.UpdatedAt = model.Time(timeutil.Now())
	a.DB.Save(&t)
	return t, nil
}

func (a *App) TenantOf(id uint64) (model.Tenant, error) {
	var t model.Tenant
	if err := a.DB.First(&t, id).Error; err != nil {
		return t, apperr.NotFound
	}
	return t, nil
}
