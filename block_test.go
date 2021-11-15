package notionapi_test

import (
	"context"
	"encoding/json"
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
			want       *notionapi.AppendBlockChildrenResponse
			wantErr    bool
			err        error
		}{
			{
				name:       "return list object",
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
				want: &notionapi.AppendBlockChildrenResponse{
					Object: notionapi.ObjectTypeList,
					Results: []notionapi.Block{
						notionapi.ParagraphBlock{
							Object:         notionapi.ObjectTypeBlock,
							ID:             "some_id",
							CreatedTime:    &timestamp,
							LastEditedTime: &timestamp,
							Type:           notionapi.BlockTypeParagraph,
							Paragraph: notionapi.Paragraph{
								Text: []notionapi.RichText{
									{
										Type: notionapi.ObjectTypeText,
										Text: notionapi.Text{Content: "AAAAAA"},
										Annotations: &notionapi.Annotations{
											Bold:  true,
											Color: notionapi.ColorDefault,
										},
										PlainText: "AAAAAA",
									},
								},
							},
						},
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
				a, err := json.Marshal(got)
				if err != nil {
					t.Errorf("AppendChidlren() marhall error = %v", err)
					return
				}
				b, err := json.Marshal(tt.want)
				if err != nil {
					t.Errorf("AppendChidlren() marhall error = %v", err)
					return
				}

				if !(string(a) == string(b)) {
					t.Errorf("AppendChidlren() got = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("Get", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.BlockID
			want       notionapi.Block
			wantErr    bool
			err        error
		}{
			{
				name:       "returns block object",
				filePath:   "testdata/block_get.json",
				statusCode: http.StatusOK,
				id:         "some_id",
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
				wantErr: false,
				err:     nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))
				got, err := client.Block.Get(context.Background(), tt.id)

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

	t.Run("Update", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.BlockID
			req        *notionapi.BlockUpdateRequest
			want       notionapi.Block
			wantErr    bool
			err        error
		}{
			{
				name:       "updates block and returns it",
				filePath:   "testdata/block_update.json",
				statusCode: http.StatusOK,
				id:         "some_id",
				req: &notionapi.BlockUpdateRequest{
					Paragraph: &notionapi.Paragraph{
						Text: []notionapi.RichText{
							{
								Text: notionapi.Text{Content: "Hello"},
							},
						},
					},
				},
				want: &notionapi.ParagraphBlock{
					Object:         notionapi.ObjectTypeBlock,
					ID:             "some_id",
					Type:           notionapi.BlockTypeParagraph,
					CreatedTime:    &timestamp,
					LastEditedTime: &timestamp,
					Paragraph: notionapi.Paragraph{
						Text: []notionapi.RichText{
							{
								Type: notionapi.ObjectTypeText,
								Text: notionapi.Text{
									Content: "Hello",
								},
								Annotations: &notionapi.Annotations{Color: notionapi.ColorDefault},
								PlainText:   "Hello",
							},
						},
					},
				},
				wantErr: false,
				err:     nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))
				got, err := client.Block.Update(context.Background(), tt.id, tt.req)

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
