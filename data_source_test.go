package notionapi_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/jomei/notionapi"
)

func TestDataSourceClient(t *testing.T) {
	timestamp, err := time.Parse(time.RFC3339, "2021-05-24T05:06:34.827Z")
	if err != nil {
		t.Fatal(err)
	}

	var user = notionapi.User{
		Object: "user",
		ID:     "some_id",
	}

	t.Run("Query", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.DataSourceID
			request    *notionapi.DataSourceQueryRequest
			want       *notionapi.DataSourceQueryResponse
			wantErr    bool
			err        error
		}{
			{
				name:       "returns query results",
				id:         "some_id",
				filePath:   "testdata/data_source_query.json",
				statusCode: http.StatusOK,
				request: &notionapi.DataSourceQueryRequest{
					Filter: &notionapi.PropertyFilter{
						Property: "Name",
						RichText: &notionapi.TextFilterCondition{
							Contains: "Hel",
						},
					},
				},
				want: &notionapi.DataSourceQueryResponse{
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
								Type:       notionapi.ParentTypeDataSourceID,
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
				got, err := client.DataSource.Query(context.Background(), tt.id, tt.request)

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

func TestDataSourceQueryRequest_MarshalJSON(t *testing.T) {
	timeObj, err := time.Parse(time.RFC3339, "2021-05-10T02:43:42Z")
	if err != nil {
		t.Error(err)
		return
	}
	dateObj := notionapi.Date(timeObj)
	tests := []struct {
		name    string
		req     *notionapi.DataSourceQueryRequest
		want    []byte
		wantErr bool
	}{
		{
			name: "timestamp created",
			req: &notionapi.DataSourceQueryRequest{
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
			req: &notionapi.DataSourceQueryRequest{
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
			req: &notionapi.DataSourceQueryRequest{
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
			req: &notionapi.DataSourceQueryRequest{
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
			req: &notionapi.DataSourceQueryRequest{
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
