package models

import "fmt"

type DateObject struct {
	Start    *Date `json:"start"`
	End      *Date `json:"end"`
	DateOnly bool  `json:"date_only,omitempty"`
}

// NewDateObject creates a new DateObject with validation
func NewDateObject(start, end *Date, dateOnly bool) (*DateObject, error) {
	// Validation: end date should not be before start date
	if start != nil && end != nil && end.Time.Before(start.Time) {
		return nil, fmt.Errorf("end date cannot be before start date")
	}

	if start != nil {
		start.DateOnly = dateOnly
	}
	if end != nil {
		end.DateOnly = dateOnly
	}

	return &DateObject{
		Start:    start,
		End:      end,
		DateOnly: dateOnly,
	}, nil
}

// FormatForNotion formats both dates with error handling
func (do DateObject) FormatForNotion() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if do.Start == nil {
		return nil, fmt.Errorf("start date is required")
	}

	result["start"] = do.Start.FormatForNotion()
	if do.End != nil {
		result["end"] = do.End.FormatForNotion()
	}

	return result, nil
}
