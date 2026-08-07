package models

import (
	"encoding/json"
	"time"
)

type Date struct {
	time.Time
	DateOnly bool
}

// Format returns the date in the specified format based on DateOnly flag
func (d Date) Format() string {
	if d.DateOnly {
		return d.Time.Format("2006-01-02")
	}
	return d.Time.Format(time.RFC3339)
}

// FormatForNotion returns the date in Notion's expected format
func (d Date) FormatForNotion() string {
	if d.DateOnly {
		return d.Time.Format("2006-01-02")
	}
	// Notion expects RFC3339 format without explicit timezone
	return d.Time.UTC().Format("2006-01-02T15:04:05Z")
}

// MarshalJSON adds error handling and validation
func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(d.FormatForNotion())
}

// UnmarshalJSON adds better error handling
func (d *Date) UnmarshalJSON(data []byte) error {
	// Handle null value
	if string(data) == "null" {
		d.Time = time.Time{}
		d.DateOnly = false
		return nil
	}

	var rawValue string
	if err := json.Unmarshal(data, &rawValue); err != nil {
		return err
	}

	// Try both formats
	formats := []string{
		"2006-01-02",           // Date only
		"2006-01-02T15:04:05Z", // Full datetime
		time.RFC3339,           // Fallback
	}

	var lastErr error
	for _, format := range formats {
		if t, err := time.Parse(format, rawValue); err == nil {
			d.Time = t
			d.DateOnly = format == "2006-01-02"
			return nil
		} else {
			lastErr = err
		}
	}

	return lastErr
}
