package timeutil

import "time"

// Beijing is GMT+8. All business timestamps must use this zone.
var Beijing = time.FixedZone("CST", 8*3600)

func Now() time.Time {
	return time.Now().In(Beijing)
}

func InBeijing(t time.Time) time.Time {
	return t.In(Beijing)
}

func ParseRFC3339(s string) (time.Time, error) {
	if s == "" {
		return Now(), nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t, err = time.Parse("2006-01-02 15:04:05", s)
		if err != nil {
			return time.Time{}, err
		}
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, Beijing)
	}
	return t.In(Beijing), nil
}
