package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/miniqxt/backend/internal/config"
	"github.com/miniqxt/backend/internal/model"
	"github.com/miniqxt/backend/internal/pkg/jwtutil"
	"github.com/miniqxt/backend/internal/router"
	"github.com/miniqxt/backend/internal/service"
)

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func TestValidationMessagesAreRequestLocal(t *testing.T) {
	const secret = "request-local-errors-test-secret"
	token, err := jwtutil.Sign(secret, 42, 7, model.RoleTenant, "Test Admin")
	if err != nil {
		t.Fatal(err)
	}

	h := router.New(&service.App{Cfg: config.Config{JWTSecret: secret}})

	status, first := post(t, h, token, "/api/v1/exam-sessions/9/traces", `{"events":[]}`)
	if status != http.StatusBadRequest || first.Code != "VALIDATION" || first.Message != "events 不能为空" {
		t.Fatalf("empty trace events: status=%d response=%+v", status, first)
	}

	status, second := post(t, h, token, "/api/v1/employees", `{`)
	if status != http.StatusBadRequest || second.Code != "VALIDATION" || second.Message != "参数校验失败" {
		t.Fatalf("malformed employee request: status=%d response=%+v", status, second)
	}
}

func post(t *testing.T, h http.Handler, token, path, body string) (int, errorResponse) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	var out errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode %s response (status %d): %v; body=%q", path, rec.Code, err, rec.Body.String())
	}
	return rec.Code, out
}
