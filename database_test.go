package notionapi_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/jomei/notionapi"
)

func TestDatabaseClient(t *testing.T) {
	timestamp, err := time.Parse(time.RFC3339, "2021-05-24T05:06:34.827Z")
	if err != nil {
		t.Fatal(err)
	}

	emoji := notionapi.Emoji("🎉")

	var user = notionapi.User{
		Object: "user",
		ID:     "some_id",
	}

	t.Run("Get", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.DatabaseID
			want       *notionapi.Database
			wantErr    bool
			err        error
		}{
			{
				name:       "returns database by id",
				id:         "some_id",
				filePath:   "testdata/database_get.json",
				statusCode: http.StatusOK,
				want: &notionapi.Database{
					Object:         notionapi.ObjectTypeDatabase,
					ID:             "some_id",
					CreatedTime:    timestamp,
					LastEditedTime: timestamp,
					CreatedBy:      user,
					LastEditedBy:   user,
					Title: []notionapi.RichText{
						{
							Type:        notionapi.RichTextTypeText,
							Text:        &notionapi.Text{Content: "Test Database"},
							Annotations: &notionapi.Annotations{Color: "default"},
							PlainText:   "Test Database",
							Href:        "",
						},
					},
					// Properties: notionapi.PropertyConfigs{
					//	"Tags": notionapi.MultiSelectPropertyConfig{
					//		ID:          ";s|V",
					//		Type:        notionapi.PropertyConfigTypeMultiSelect,
					//		MultiSelect: notionapi.Select{Options: []notionapi.Option{{ID: "id", Name: "tag", Color: "Blue"}}},
					//	},
					//	"Some another column": notionapi.PeoplePropertyConfig{
					//		ID:   "rJt\\",
					//		Type: notionapi.PropertyConfigTypePeople,
					//	},
					//
					//	"Name": notionapi.TitlePropertyConfig{
					//		ID:    "title",
					//		Type:  notionapi.PropertyConfigTypeTitle,
					//		Title: notionapi.RichText{},
					//	},
					// },
				},
				wantErr: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))
				got, err := client.Database.Get(context.Background(), tt.id)

				if (err != nil) != tt.wantErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				// TODO: remove properties from comparing for a while. Have to compare with interface somehow
				got.Properties = nil
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("Update", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.DatabaseID
			request    *notionapi.DatabaseUpdateRequest
			want       *notionapi.Database
			wantErr    bool
			err        error
		}{
			{
				name:       "returns update results",
				filePath:   "testdata/database_update.json",
				statusCode: http.StatusOK,
				id:         "some_id",
				request: &notionapi.DatabaseUpdateRequest{
					Title: []notionapi.RichText{
						{
							Type: notionapi.RichTextTypeText,
							Text: &notionapi.Text{Content: "patch"},
						},
					},
					Properties: notionapi.PropertyConfigs{
						"patch": notionapi.TitlePropertyConfig{
							Type: notionapi.PropertyConfigTypeRichText,
						},
					},
				},
				want: &notionapi.Database{
					Object:         notionapi.ObjectTypeDatabase,
					ID:             "some_id",
					CreatedTime:    timestamp,
					LastEditedTime: timestamp,
					CreatedBy:      user,
					LastEditedBy:   user,
					Parent: notionapi.Parent{
						Type:   "page_id",
						PageID: "48f8fee9-cd79-4180-bc2f-ec0398253067",
					},
					Title: []notionapi.RichText{
						{
							Type: notionapi.RichTextTypeText,
							Text: &notionapi.Text{Content: "patch"},
						},
					},
					Description: []notionapi.RichText{},
					IsInline:    false,
					Archived:    false,
					Icon: &notionapi.Icon{
						Type:  "emoji",
						Emoji: &emoji,
					},
					Cover: &notionapi.Image{
						Type: "external",
						External: &notionapi.FileObject{
							URL: "https://website.domain/images/image.png",
						},
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))
				got, err := client.Database.Update(context.Background(), tt.id, tt.request)

				if (err != nil) != tt.wantErr {
					t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				got.Properties = nil
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Update() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("Create", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			request    *notionapi.DatabaseCreateRequest
			want       *notionapi.Database
			wantErr    bool
			err        error
		}{
			{
				name:       "returns created db",
				filePath:   "testdata/database_create.json",
				statusCode: http.StatusOK,
				request: &notionapi.DatabaseCreateRequest{
					Parent: notionapi.Parent{
						Type:   notionapi.ParentTypePageID,
						PageID: "some_id",
					},
					Title: []notionapi.RichText{
						{
							Type: notionapi.RichTextTypeText,
							Text: &notionapi.Text{Content: "Grocery List"},
						},
					},
					Properties: notionapi.PropertyConfigs{
						"create": notionapi.TitlePropertyConfig{
							Type: notionapi.PropertyConfigTypeTitle,
						},
					},
					IsInline: false,
				},
				want: &notionapi.Database{
					Object:         notionapi.ObjectTypeDatabase,
					ID:             "some_id",
					CreatedTime:    timestamp,
					LastEditedTime: timestamp,
					CreatedBy:      user,
					LastEditedBy:   user,
					Parent: notionapi.Parent{
						Type:   "page_id",
						PageID: "a7744006-9233-4cd0-bf44-3a49de2c01b5",
					},
					Title: []notionapi.RichText{
						{
							Type:        notionapi.RichTextTypeText,
							Text:        &notionapi.Text{Content: "Grocery List"},
							PlainText:   "Grocery List",
							Annotations: &notionapi.Annotations{Color: notionapi.ColorDefault},
						},
					},
					Description: []notionapi.RichText{},
					IsInline:    false,
					Archived:    false,
					Icon: &notionapi.Icon{
						Type:  "emoji",
						Emoji: &emoji,
					},
					Cover: &notionapi.Image{
						Type: "external",
						External: &notionapi.FileObject{
							URL: "https://website.domain/images/image.png",
						},
					},
				},
			},
			{
				name:       "returns created db 2",
				filePath:   "testdata/database_create_2.json",
				statusCode: http.StatusOK,
				request: &notionapi.DatabaseCreateRequest{
					Parent: notionapi.Parent{
						Type:   notionapi.ParentTypePageID,
						PageID: "some_id",
					},
					Title: []notionapi.RichText{
						{
							Type: notionapi.RichTextTypeText,
							Text: &notionapi.Text{Content: "Grocery List"},
						},
					},
					Properties: notionapi.PropertyConfigs{
						"create": notionapi.TitlePropertyConfig{
							Type: notionapi.PropertyConfigTypeTitle,
						},
					},
					IsInline: false,
				},
				want: &notionapi.Database{
					Object:         notionapi.ObjectTypeDatabase,
					ID:             "some_id",
					CreatedTime:    timestamp,
					LastEditedTime: timestamp,
					CreatedBy:      user,
					LastEditedBy:   user,
					Parent: notionapi.Parent{
						Type:    "block_id",
						BlockID: "a7744006-9233-4cd0-bf44-3a49de2c01b5",
					},
					Title: []notionapi.RichText{
						{
							Type:        notionapi.RichTextTypeText,
							Text:        &notionapi.Text{Content: "Grocery List"},
							PlainText:   "Grocery List",
							Annotations: &notionapi.Annotations{Color: notionapi.ColorDefault},
						},
					},
					Description: []notionapi.RichText{},
					IsInline:    false,
					Archived:    false,
					Icon: &notionapi.Icon{
						Type:  "emoji",
						Emoji: &emoji,
					},
					Cover: &notionapi.Image{
						Type: "external",
						External: &notionapi.FileObject{
							URL: "https://website.domain/images/image.png",
						},
					},
				},
			},
			{
				name:       "returns created db 3",
				filePath:   "testdata/database_create_3.json",
				statusCode: http.StatusOK,
				request: &notionapi.DatabaseCreateRequest{
					Parent: notionapi.Parent{
						Type:   notionapi.ParentTypePageID,
						PageID: "some_id",
					},
					Title: []notionapi.RichText{
						{
							Type: notionapi.RichTextTypeText,
							Text: &notionapi.Text{Content: "Grocery List"},
						},
					},
					Properties: notionapi.PropertyConfigs{
						"create": notionapi.TitlePropertyConfig{
							Type: notionapi.PropertyConfigTypeTitle,
						},
					},
					IsInline: true,
				},
				want: &notionapi.Database{
					Object:         notionapi.ObjectTypeDatabase,
					ID:             "some_id",
					CreatedTime:    timestamp,
					LastEditedTime: timestamp,
					CreatedBy:      user,
					LastEditedBy:   user,
					Parent: notionapi.Parent{
						Type:   "page_id",
						PageID: "a7744006-9233-4cd0-bf44-3a49de2c01b5",
					},
					Title: []notionapi.RichText{
						{
							Type:        notionapi.RichTextTypeText,
							Text:        &notionapi.Text{Content: "Grocery List"},
							PlainText:   "Grocery List",
							Annotations: &notionapi.Annotations{Color: notionapi.ColorDefault},
						},
					},
					Description: []notionapi.RichText{},
					IsInline:    true,
					Archived:    false,
					Icon: &notionapi.Icon{
						Type:  "emoji",
						Emoji: &emoji,
					},
					Cover: &notionapi.Image{
						Type: "external",
						External: &notionapi.FileObject{
							URL: "https://website.domain/images/image.png",
						},
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))

				got, err := client.Database.Create(context.Background(), tt.request)

				if (err != nil) != tt.wantErr {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				got.Properties = nil
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Create() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("Get with empty database_id", func(t *testing.T) {
		client := notionapi.NewClient("some_token")
		_, err := client.Database.Get(context.TODO(), notionapi.DatabaseID(""))
		if err.Error() != "empty database id" {
			t.Error("database id is required error is expected")
		}
	})
}
