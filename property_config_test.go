package notionapi

import (
	"reflect"
	"testing"
)

func TestNumberPropertyConfig_MarshalJSON(t *testing.T) {
	type fields struct {
		ID     ObjectID
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
			name: "The Format field goes into the number property",
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
				ID:     tt.fields.ID,
				Type:   tt.fields.Type,
				Format: tt.fields.Format,
			}
			got, err := p.MarshalJSON()
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
