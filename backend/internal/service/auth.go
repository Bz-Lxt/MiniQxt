package service

import (
	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/pkg/apperr"
	"github.com/miniqxt/backend/internal/pkg/jwtutil"
	"github.com/miniqxt/backend/internal/pkg/passwd"
)

func (a *App) Login(username, password string) (string, model.PublicUser, error) {
	if username == "" || password == "" {
		return "", model.PublicUser{}, apperr.Validation.With("用户名和密码不能为空")
	}
	var u model.User
	if err := a.DB.Where("username = ?", username).First(&u).Error; err != nil {
		return "", model.PublicUser{}, apperr.Unauthorized.With("用户名或密码错误")
	}
	if u.Status != model.StatusActive {
		return "", model.PublicUser{}, apperr.Forbidden.With("账号已停用")
	}
	if u.TenantID > 0 {
		var t model.Tenant
		if err := a.DB.First(&t, u.TenantID).Error; err != nil || t.Status != model.StatusActive {
			return "", model.PublicUser{}, apperr.Forbidden.With("租户已停用")
		}
	}
	if !passwd.Verify(u.PasswordHash, password) {
		return "", model.PublicUser{}, apperr.Unauthorized.With("用户名或密码错误")
	}
	tok, err := jwtutil.Sign(a.Cfg.JWTSecret, u.ID, u.TenantID, u.Role, u.DisplayName)
	if err != nil {
		return "", model.PublicUser{}, apperr.Internal
	}
	return tok, u.Public(), nil
}

func (a *App) Me(id uint64) (model.User, error) {
	var u model.User
	if err := a.DB.Preload("Department").Preload("Position").First(&u, id).Error; err != nil {
		return u, apperr.NotFound
	}
	return u, nil
}
