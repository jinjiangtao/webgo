package utils

import (
	"database/sql/driver"
	"errors"
	"time"
)

type DateTime struct {
	time.Time
}

func (d *DateTime) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		d.Time = time.Time{}
		return nil
	}

	s := string(data)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	tzFormats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
	}

	for _, format := range tzFormats {
		t, err := time.Parse(format, s)
		if err == nil {
			d.Time = t
			return nil
		}
	}

	localFormats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006-01-02T15:04:05.000",
		"2006-01-02 15:04:05.000",
	}

	for _, format := range localFormats {
		t, err := time.ParseInLocation(format, s, time.Local)
		if err == nil {
			d.Time = t
			return nil
		}
	}

	return errors.New("invalid time format: " + s)
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Time.Format(time.RFC3339) + `"`), nil
}

func (d DateTime) Value() (driver.Value, error) {
	return d.Time, nil
}

func (d *DateTime) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return errors.New("invalid time type")
	}
	d.Time = t
	return nil
}

func (d DateTime) ToTime() time.Time {
	return d.Time
}
