package notionapi_test

import (
	"context"
	"github.com/jomei/notionapi"
	"reflect"
	"testing"
	"time"
)

func TestDatabaseClient(t *testing.T) {
	timestamp, err := time.Parse(time.RFC3339, "2021-05-24T05:06:34.827Z")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("#Get", func(t *testing.T) {
		tests := []struct {
			name     string
			filePath string
			id       notionapi.DatabaseID
			want     *notionapi.DatabaseObject
			wantErr  bool
			err      error
		}{
			{
				name: "returns database by id",
				id:   "some_id",
				want: &notionapi.DatabaseObject{
					Object:         notionapi.ObjectTypeDatabase,
					ID:             "some_id",
					CreatedTime:    timestamp,
					LastEditedTime: timestamp,
					Title: []notionapi.RichTextObject{
						{
							Type:        notionapi.ObjectTypeText,
							Text:        notionapi.TextObject{Content: "Test Database", Link: ""},
							Annotations: notionapi.Annotations{Color: "default"},
							PlainText:   "Test Database",
							Href:        "",
						},
					},

					Properties: map[notionapi.PropertyName]notionapi.BasicObject{
						"Tags": {
							ID:          ";s|V",
							Type:        "multi_select",
							MultiSelect: &notionapi.MultiSelectObject{Options: []notionapi.SelectOption{{ID: "id", Name: "tag", Color: "Blue"}}},
						},
						"Some another column": {
							ID:     "rJt\\",
							Type:   notionapi.ObjectTypePeople,
							People: &struct{}{},
						},
						"SomeColumn": {
							ID:       "~j_@",
							Type:     notionapi.ObjectTypeRichText,
							RichText: &notionapi.RichTextObject{},
						},
						"Name": {
							ID:    "title",
							Type:  notionapi.ObjectTypeTitle,
							Title: &notionapi.TextObject{},
						},
					},
				},
				wantErr: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {

				c := newMockedClient(t, tt.filePath)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))
				got, err := client.Database.Get(context.Background(), tt.id)

				if (err != nil) != tt.wantErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
			})
		}
	})
}
