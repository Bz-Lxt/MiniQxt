package jwtutil

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/miniqxt/backend/internal/timeutil"
)

type Claims struct {
	UserID   uint64 `json:"uid"`
	TenantID uint64 `json:"tid"`
	Role     string `json:"role"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

func Sign(secret string, userID, tenantID uint64, role, name string) (string, error) {
	now := timeutil.Now()
	c := Claims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		Name:     name,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			Issuer:    "miniqxt",
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString([]byte(secret))
}

func Parse(secret, token string) (*Claims, error) {
	t, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	c, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return c, nil
}
