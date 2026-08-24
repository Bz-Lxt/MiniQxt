package model

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/miniqxt/backend/internal/timeutil"
)

// Time stores and serializes timestamps in GMT+8.
type Time time.Time

func (t Time) Time() time.Time { return time.Time(t) }

func (t Time) IsZero() bool { return time.Time(t).IsZero() }

func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	s := time.Time(t).In(timeutil.Beijing).Format("2006-01-02T15:04:05+08:00")
	return []byte(`"` + s + `"`), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == "null" || s == `""` {
		*t = Time{}
		return nil
	}
	if len(s) >= 2 && s[0] == '"' {
		s = s[1 : len(s)-1]
	}
	parsed, err := timeutil.ParseRFC3339(s)
	if err != nil {
		return err
	}
	*t = Time(parsed)
	return nil
}

func (t Time) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return time.Time(t).In(timeutil.Beijing), nil
}

func (t *Time) Scan(v any) error {
	if v == nil {
		*t = Time{}
		return nil
	}
	switch x := v.(type) {
	case time.Time:
		*t = Time(x.In(timeutil.Beijing))
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Time", v)
	}
}

const (
	RolePlatform = "platform_admin"
	RoleTenant   = "tenant_admin"
	RoleInstructor = "instructor"
	RoleEmployee = "employee"

	StatusActive   = "active"
	StatusDisabled = "disabled"

	QSingle = "single"
	QMulti  = "multi"
	QJudge  = "judge"
	QEssay  = "essay"

	SessNotStarted = "not_started"
	SessInProgress = "in_progress"
	SessSubmitted  = "submitted"
	SessForced     = "force_submitted"
	SessExpired    = "expired"

	SubQueued        = "queued"
	SubGrading       = "grading"
	SubGraded        = "graded"
	SubPendingManual = "pending_manual"
	SubFailed        = "failed"

	IntegrityOK         = "ok"
	IntegrityUnverified = "integrity_unverified"
	IntegritySuspicious = "suspicious"

	CertPending = "pending"
	CertIssued  = "issued"
	CertRevoked = "revoked"
)
