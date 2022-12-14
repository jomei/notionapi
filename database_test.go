package notionapi_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/conduitio-labs/notionapi"
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
							Type:        notionapi.ObjectTypeText,
							Text:        &notionapi.Text{Content: "Test Database"},
							Annotations: &notionapi.Annotations{Color: "default"},
							PlainText:   "Test Database",
							Href:        "",
						},
					},
					//Properties: notionapi.PropertyConfigs{
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
					//},
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

	t.Run("Query", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.DatabaseID
			request    *notionapi.DatabaseQueryRequest
			want       *notionapi.DatabaseQueryResponse
			wantErr    bool
			err        error
		}{
			{
				name:       "returns query results",
				id:         "some_id",
				filePath:   "testdata/database_query.json",
				statusCode: http.StatusOK,
				request: &notionapi.DatabaseQueryRequest{
					Filter: &notionapi.PropertyFilter{
						Property: "Name",
						RichText: &notionapi.TextFilterCondition{
							Contains: "Hel",
						},
					},
				},
				want: &notionapi.DatabaseQueryResponse{
					Object: notionapi.ObjectTypeList,
					Results: []notionapi.Page{
						{
							Object:         notionapi.ObjectTypePage,
							ID:             "some_id",
							CreatedTime:    timestamp,
							LastEditedTime: timestamp,
							CreatedBy:      user,
							LastEditedBy:   user,
							Parent: notionapi.Parent{
								Type:       notionapi.ParentTypeDatabaseID,
								DatabaseID: "some_id",
							},
							Archived: false,
							URL:      "some_url",
						},
					},
					HasMore:    false,
					NextCursor: "",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))
				got, err := client.Database.Query(context.Background(), tt.id, tt.request)

				if (err != nil) != tt.wantErr {
					t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				got.Results[0].Properties = nil
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Query() got = %v, want %v", got, tt.want)
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
							Type: notionapi.ObjectTypeText,
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
							Type: notionapi.ObjectTypeText,
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
							Type: notionapi.ObjectTypeText,
							Text: &notionapi.Text{Content: "Grocery List"},
						},
					},
					Properties: notionapi.PropertyConfigs{
						"create": notionapi.TitlePropertyConfig{
							Type: notionapi.PropertyConfigTypeTitle,
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
						PageID: "a7744006-9233-4cd0-bf44-3a49de2c01b5",
					},
					Title: []notionapi.RichText{
						{
							Type:        notionapi.ObjectTypeText,
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
							Type: notionapi.ObjectTypeText,
							Text: &notionapi.Text{Content: "Grocery List"},
						},
					},
					Properties: notionapi.PropertyConfigs{
						"create": notionapi.TitlePropertyConfig{
							Type: notionapi.PropertyConfigTypeTitle,
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
						Type:    "block_id",
						BlockID: "a7744006-9233-4cd0-bf44-3a49de2c01b5",
					},
					Title: []notionapi.RichText{
						{
							Type:        notionapi.ObjectTypeText,
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
}

func TestDatabaseQueryRequest_MarshalJSON(t *testing.T) {
	timeObj, err := time.Parse(time.RFC3339, "2021-05-10T02:43:42Z")
	if err != nil {
		t.Error(err)
		return
	}
	dateObj := notionapi.Date(timeObj)
	tests := []struct {
		name    string
		req     *notionapi.DatabaseQueryRequest
		want    []byte
		wantErr bool
	}{
		{
			name: "timestamp created",
			req: &notionapi.DatabaseQueryRequest{
				Filter: &notionapi.TimestampFilter{
					Timestamp: notionapi.TimestampCreated,
					CreatedTime: &notionapi.DateFilterCondition{
						NextWeek: &struct{}{},
					},
				},
			},
			want: []byte(`{"filter":{"timestamp":"created_time","created_time":{"next_week":{}}}}`),
		},
		{
			name: "timestamp last edited",
			req: &notionapi.DatabaseQueryRequest{
				Filter: &notionapi.TimestampFilter{
					Timestamp: notionapi.TimestampLastEdited,
					LastEditedTime: &notionapi.DateFilterCondition{
						Before: &dateObj,
					},
				},
			},
			want: []byte(`{"filter":{"timestamp":"last_edited_time","last_edited_time":{"before":"2021-05-10T02:43:42Z"}}}`),
		},
		{
			name: "or compound filter one level",
			req: &notionapi.DatabaseQueryRequest{
				Filter: notionapi.OrCompoundFilter{
					notionapi.PropertyFilter{
						Property: "Status",
						Select: &notionapi.SelectFilterCondition{
							Equals: "Reading",
						},
					},
					notionapi.PropertyFilter{
						Property: "Publisher",
						Select: &notionapi.SelectFilterCondition{
							Equals: "NYT",
						},
					},
				},
			},
			want: []byte(`{"filter":{"or":[{"property":"Status","select":{"equals":"Reading"}},{"property":"Publisher","select":{"equals":"NYT"}}]}}`),
		},
		{
			name: "and compound filter one level",
			req: &notionapi.DatabaseQueryRequest{
				Filter: notionapi.AndCompoundFilter{
					notionapi.PropertyFilter{
						Property: "Status",
						Select: &notionapi.SelectFilterCondition{
							Equals: "Reading",
						},
					},
					notionapi.PropertyFilter{
						Property: "Publisher",
						Select: &notionapi.SelectFilterCondition{
							Equals: "NYT",
						},
					},
				},
			},
			want: []byte(`{"filter":{"and":[{"property":"Status","select":{"equals":"Reading"}},{"property":"Publisher","select":{"equals":"NYT"}}]}}`),
		},
		{
			name: "compound filter two levels",
			req: &notionapi.DatabaseQueryRequest{
				Filter: notionapi.OrCompoundFilter{
					notionapi.PropertyFilter{
						Property: "Description",
						RichText: &notionapi.TextFilterCondition{
							Contains: "fish",
						},
					},
					notionapi.AndCompoundFilter{
						notionapi.PropertyFilter{
							Property: "Food group",
							Select: &notionapi.SelectFilterCondition{
								Equals: "🥦Vegetable",
							},
						},
						notionapi.PropertyFilter{
							Property: "Is protein rich?",
							Checkbox: &notionapi.CheckboxFilterCondition{
								Equals: true,
							},
						},
					},
				},
			},
			want: []byte(`{"filter":{"or":[{"property":"Description","rich_text":{"contains":"fish"}},{"and":[{"property":"Food group","select":{"equals":"🥦Vegetable"}},{"property":"Is protein rich?","checkbox":{"equals":true}}]}]}}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.req.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalJSON() got = %s, want %s", got, tt.want)
			}
		})
	}
}
