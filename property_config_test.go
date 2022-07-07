package notionapi

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNumberPropertyConfig_MarshalJSON(t *testing.T) {
	type fields struct {
		Type   PropertyConfigType
		Format FormatType
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "returns correct json",
			fields: fields{
				Type:   PropertyConfigTypeNumber,
				Format: FormatDollar,
			},
			want:    []byte(`{"type":"number","number":{"format":"dollar"}}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NumberPropertyConfig{
				Type:   tt.fields.Type,
				Number: NumberFormat{Format: tt.fields.Format},
			}
			got, err := json.Marshal(p)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
