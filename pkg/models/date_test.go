package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDate_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		date    Date
		want    string
		wantErr bool
	}{
		{
			name: "date only format",
			date: Date{
				Time:     time.Date(2024, 3, 14, 0, 0, 0, 0, time.UTC),
				DateOnly: true,
			},
			want:    `"2024-03-14"`,
			wantErr: false,
		},
		{
			name: "datetime format",
			date: Date{
				Time:     time.Date(2024, 3, 14, 15, 30, 0, 0, time.UTC),
				DateOnly: false,
			},
			want:    `"2024-03-14T15:30:00Z"`,
			wantErr: false,
		},
		{
			name: "zero time",
			date: Date{
				Time:     time.Time{},
				DateOnly: false,
			},
			want:    "null",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Date.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    Date
		wantErr bool
	}{
		{
			name: "date only format",
			json: `"2024-03-14"`,
			want: Date{
				Time:     time.Date(2024, 3, 14, 0, 0, 0, 0, time.UTC),
				DateOnly: true,
			},
			wantErr: false,
		},
		{
			name: "datetime format",
			json: `"2024-03-14T15:30:00Z"`,
			want: Date{
				Time:     time.Date(2024, 3, 14, 15, 30, 0, 0, time.UTC),
				DateOnly: false,
			},
			wantErr: false,
		},
		{
			name:    "invalid format",
			json:    `"invalid-date"`,
			wantErr: true,
		},
		{
			name: "null value",
			json: "null",
			want: Date{
				Time:     time.Time{},
				DateOnly: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Date
			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Date.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !got.Time.Equal(tt.want.Time) {
					t.Errorf("Date.UnmarshalJSON() Time = %v, want %v", got.Time, tt.want.Time)
				}
				if got.DateOnly != tt.want.DateOnly {
					t.Errorf("Date.UnmarshalJSON() DateOnly = %v, want %v", got.DateOnly, tt.want.DateOnly)
				}
			}
		})
	}
}

func TestDateObject_FormatForNotion(t *testing.T) {
	tests := []struct {
		name    string
		obj     DateObject
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid date range",
			obj: DateObject{
				Start: &Date{
					Time:     time.Date(2024, 3, 14, 0, 0, 0, 0, time.UTC),
					DateOnly: true,
				},
				End: &Date{
					Time:     time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
					DateOnly: true,
				},
				DateOnly: true,
			},
			want: map[string]interface{}{
				"start": "2024-03-14",
				"end":   "2024-03-15",
			},
			wantErr: false,
		},
		{
			name: "start date only",
			obj: DateObject{
				Start: &Date{
					Time:     time.Date(2024, 3, 14, 15, 30, 0, 0, time.UTC),
					DateOnly: false,
				},
				DateOnly: false,
			},
			want: map[string]interface{}{
				"start": "2024-03-14T15:30:00Z",
			},
			wantErr: false,
		},
		{
			name: "missing start date",
			obj: DateObject{
				End: &Date{
					Time:     time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
					DateOnly: true,
				},
				DateOnly: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.obj.FormatForNotion()
			if (err != nil) != tt.wantErr {
				t.Errorf("DateObject.FormatForNotion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				for k, v := range tt.want {
					if got[k] != v {
						t.Errorf("DateObject.FormatForNotion() = %v, want %v", got[k], v)
					}
				}
			}
		})
	}
}

func TestNewDateObject(t *testing.T) {
	tests := []struct {
		name     string
		start    *Date
		end      *Date
		dateOnly bool
		wantErr  bool
	}{
		{
			name: "valid date range",
			start: &Date{
				Time:     time.Date(2024, 3, 14, 0, 0, 0, 0, time.UTC),
				DateOnly: true,
			},
			end: &Date{
				Time:     time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
				DateOnly: true,
			},
			dateOnly: true,
			wantErr:  false,
		},
		{
			name: "end before start",
			start: &Date{
				Time:     time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
				DateOnly: true,
			},
			end: &Date{
				Time:     time.Date(2024, 3, 14, 0, 0, 0, 0, time.UTC),
				DateOnly: true,
			},
			dateOnly: true,
			wantErr:  true,
		},
		{
			name:     "start date only",
			start:    &Date{Time: time.Now()},
			end:      nil,
			dateOnly: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDateObject(tt.start, tt.end, tt.dateOnly)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDateObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
