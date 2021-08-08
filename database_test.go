package notionapi_test

import (
	"context"
	"github.com/jomei/notionapi"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestDatabaseClient(t *testing.T) {
	timestamp, err := time.Parse(time.RFC3339, "2021-05-24T05:06:34.827Z")
	if err != nil {
		t.Fatal(err)
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
					Title: []notionapi.RichText{
						{
							Type:        notionapi.ObjectTypeText,
							Text:        notionapi.Text{Content: "Test Database", Link: ""},
							Annotations: &notionapi.Annotations{Color: "default"},
							PlainText:   "Test Database",
							Href:        "",
						},
					},

					//	Properties: notionapi.Properties{
					//		"Tags": notionapi.MultiSelectProperty{
					//			ID:          ";s|V",
					//			Type:        notionapi.PropertyTypeMultiSelect,
					//			MultiSelect: notionapi.Select{Options: []notionapi.Option{{ID: "id", Name: "tag", Color: "Blue"}}},
					//		},
					//		"Some another column": notionapi.PeopleProperty{
					//			ID:     "rJt\\",
					//			Type:   notionapi.PropertyTypePeople,
					//		},
					//		"SomeColumn": notionapi.RichTextProperty{
					//			ID:       "~j_@",
					//			Type:     notionapi.PropertyTypeRichText,
					//			RichText: notionapi.RichText{},
					//		},
					//		"Name": notionapi.TitleProperty{
					//			ID:    "title",
					//			Type:  notionapi.PropertyTypeTitle,
					//			Title: notionapi.RichText{},
					//		},
					//	},
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

	t.Run("List", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			want       *notionapi.DatabaseListResponse
			wantErr    bool
			err        error
		}{
			{
				name:       "returns list of databases",
				filePath:   "testdata/database_list.json",
				statusCode: http.StatusOK,
				want: &notionapi.DatabaseListResponse{
					Object: notionapi.ObjectTypeList,
					Results: []notionapi.Database{
						{
							Object:         notionapi.ObjectTypeDatabase,
							ID:             "some_id",
							CreatedTime:    timestamp,
							LastEditedTime: timestamp,
							Title: notionapi.Paragraph{
								{
									Type: notionapi.ObjectTypeText,
									Text: notionapi.Text{
										Content: "Test Database",
									},
									Annotations: &notionapi.Annotations{
										Color: notionapi.ColorDefault,
									},
									PlainText: "Test Database",
								},
							},
						},
					},
					HasMore: false,
				},
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))

				got, err := client.Database.List(context.Background(), nil)

				if (err != nil) != tt.wantErr {
					t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				got.Results[0].Properties = nil
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("List() got = %v, want %v", got, tt.want)
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
					PropertyFilter: &notionapi.PropertyFilter{
						Property: "Name",
						Text: &notionapi.TextFilterCondition{
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
			name: "with property filter without sort",
			req: &notionapi.DatabaseQueryRequest{
				PropertyFilter: &notionapi.PropertyFilter{
					Property: "Status",
					Select: &notionapi.SelectFilterCondition{
						Equals: "Reading",
					},
				},
			},
			want: []byte(`{"filter":{"property":"Status","select":{"equals":"Reading"}}}`),
		},
		{
			name: "with property filter with sort",
			req: &notionapi.DatabaseQueryRequest{
				PropertyFilter: &notionapi.PropertyFilter{
					Property: "Status",
					Select: &notionapi.SelectFilterCondition{
						Equals: "Reading",
					},
				},
				Sorts: []notionapi.SortObject{
					{
						Property:  "Score /5",
						Direction: notionapi.SortOrderASC,
					},
				},
			},
			want: []byte(`{"sorts":[{"property":"Score /5","direction":"ascending"}],"filter":{"property":"Status","select":{"equals":"Reading"}}}`),
		},
		{
			name: "compound filter",
			req: &notionapi.DatabaseQueryRequest{
				CompoundFilter: &notionapi.CompoundFilter{
					notionapi.FilterOperatorOR: []notionapi.PropertyFilter{
						{
							Property: "Status",
							Select: &notionapi.SelectFilterCondition{
								Equals: "Reading",
							},
						},
						{
							Property: "Publisher",
							Select: &notionapi.SelectFilterCondition{
								Equals: "NYT",
							},
						},
					},
				},
			},
			want: []byte(`{"filter":{"or":[{"property":"Status","select":{"equals":"Reading"}},{"property":"Publisher","select":{"equals":"NYT"}}]}}`),
		},
		{
			name: "date filter",
			req: &notionapi.DatabaseQueryRequest{
				PropertyFilter: &notionapi.PropertyFilter{
					Property: "created_at",
					Date: &notionapi.DateFilterCondition{
						Equals:   &dateObj,
						PastWeek: &struct{}{},
					},
				},
			},
			want: []byte(`{"filter":{"property":"created_at","date":{"equals":"2021-05-10T02:43:42Z","past_week":{}}}}`),
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
