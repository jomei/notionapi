package notionapi_test

import (
	"context"
	"github.com/jomei/notionapi"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestBlockClient(t *testing.T) {
	timestamp, err := time.Parse(time.RFC3339, "2021-05-24T05:06:34.827Z")
	if err != nil {
		t.Fatal(err)
	}
	t.Run("GetChildren", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.BlockID
			len        int
			wantErr    bool
			err        error
		}{
			{
				name:       "returns blocks by id of parent block",
				id:         "some_id",
				statusCode: http.StatusOK,
				filePath:   "testdata/block_get_children.json",
				len:        2,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))
				got, err := client.Block.GetChildren(context.Background(), tt.id, nil)

				if (err != nil) != tt.wantErr {
					t.Errorf("GetChildren() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if tt.len != len(got.Results) {
					t.Errorf("GetChildren got %d, want: %d", len(got.Results), tt.len)
				}
			})
		}
	})

	t.Run("AppendChildren", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.BlockID
			request    *notionapi.AppendBlockChildrenRequest
			want       *notionapi.ChildPageBlock
			wantErr    bool
			err        error
		}{
			{
				name:       "returns blocks by id of parent block",
				id:         "some_id",
				filePath:   "testdata/block_append_children.json",
				statusCode: http.StatusOK,
				request: &notionapi.AppendBlockChildrenRequest{
					Children: []notionapi.Block{
						&notionapi.Heading2Block{
							Object: notionapi.ObjectTypeBlock,
							Type:   notionapi.BlockTypeHeading2,
							Heading2: struct {
								Text []notionapi.RichText `json:"text"`
							}{[]notionapi.RichText{
								{
									Type: notionapi.ObjectTypeText,
									Text: notionapi.Text{Content: "Hello"},
								},
							}},
						},
					},
				},
				want: &notionapi.ChildPageBlock{
					Object:         notionapi.ObjectTypeBlock,
					ID:             "some_id",
					Type:           notionapi.BlockTypeChildPage,
					CreatedTime:    &timestamp,
					LastEditedTime: &timestamp,
					HasChildren:    true,
					ChildPage: struct {
						Title string `json:"title"`
					}{
						Title: "Hello",
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))
				got, err := client.Block.AppendChildren(context.Background(), tt.id, tt.request)

				if (err != nil) != tt.wantErr {
					t.Errorf("AppendChidlren() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("AppendChidlren() got = %v, want %v", got, tt.want)
				}
			})
		}
	})
}
