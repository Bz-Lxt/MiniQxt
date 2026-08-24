package router_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/pkg/jwtutil"
	"github.com/miniqxt/backend/internal/router"
	"github.com/miniqxt/backend/internal/service"
)

func TestInstructorCannotCreateTenant(t *testing.T) {
	const secret = "rbac-isolation-test-secret"
	token, err := jwtutil.Sign(secret, 17, 3, model.RoleInstructor, "instructor")
	if err != nil {
		t.Fatal(err)
	}

	app := &service.App{}
	app.Cfg.Env = "production"
	app.Cfg.JWTSecret = secret
	handler := router.New(app)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", bytes.NewBufferString(`{}`))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("instructor request reached platform endpoint: status=%d body=%s", rec.Code, rec.Body.String())
	}
	var body struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Code != "FORBIDDEN" {
		t.Fatalf("unexpected error code %q", body.Code)
	}
}
