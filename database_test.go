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
			want     *notionapi.Database
			wantErr  bool
			err      error
		}{
			{
				name: "returns database by id",
				id:   "some_id",
				want: &notionapi.Database{
					Object:         notionapi.ObjectTypeDatabase,
					ID:             "some_id",
					CreatedTime:    timestamp,
					LastEditedTime: timestamp,
					Title: []notionapi.TextObject{
						{
							Type:        notionapi.ObjectTypeText,
							Text:        notionapi.Text{Content: "Test Database", Link: ""},
							Annotations: notionapi.Annotations{Color: "default"},
							PlainText:   "Test Database",
							Href:        "",
						},
					},

					Properties: map[string]notionapi.Property{
						"Tags": notionapi.MultiSelectProperty{
							ID:          ";s|V",
							Type:        notionapi.PropertyTypeMultiSelect,
							MultiSelect: notionapi.Select{Options: []notionapi.Option{{ID: "id", Name: "tag", Color: "Blue"}}},
						},
						"Some another column": notionapi.PeopleProperty{
							ID:     "rJt\\",
							Type:   notionapi.PropertyTypePeople,
							People: &struct{}{},
						},
						"SomeColumn": notionapi.RichTextProperty{
							ID:       "~j_@",
							Type:     notionapi.PropertyTypeRichText,
							RichText: notionapi.TextObject{},
						},
						"Name": notionapi.TitleProperty{
							ID:    "title",
							Type:  notionapi.PropertyTypeTitle,
							Title: notionapi.TextObject{},
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
