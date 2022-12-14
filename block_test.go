package notionapi_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/conduitio-labs/notionapi"
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
							BasicBlock: notionapi.BasicBlock{
								Object: notionapi.ObjectTypeBlock,
								Type:   notionapi.BlockTypeHeading2,
							},
							Heading2: struct {
								RichText []notionapi.RichText `json:"rich_text"`
								Children notionapi.Blocks     `json:"children,omitempty"`
								Color    string               `json:"color,omitempty"`
							}{[]notionapi.RichText{
								{
									Type: notionapi.ObjectTypeText,
									Text: &notionapi.Text{Content: "Hello"},
								},
							}, nil, "",
							},
						},
					},
				},
				want: &notionapi.AppendBlockChildrenResponse{
					Object: notionapi.ObjectTypeList,
					Results: []notionapi.Block{
						notionapi.ParagraphBlock{
							BasicBlock: notionapi.BasicBlock{
								Object:         notionapi.ObjectTypeBlock,
								ID:             "some_id",
								CreatedTime:    &timestamp,
								LastEditedTime: &timestamp,
								Type:           notionapi.BlockTypeParagraph,
								CreatedBy: &notionapi.User{
									Object: "user",
									ID:     "some_id",
								},
								LastEditedBy: &notionapi.User{
									Object: "user",
									ID:     "some_id",
								},
							},
							Paragraph: notionapi.Paragraph{
								RichText: []notionapi.RichText{
									{
										Type: notionapi.ObjectTypeText,
										Text: &notionapi.Text{Content: "AAAAAA"},
										Annotations: &notionapi.Annotations{
											Bold:  true,
											Color: notionapi.ColorDefault,
										},
										PlainText: "AAAAAA",
									},
								},
								Color: "blue",
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
					t.Errorf("AppendChildren() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				a, err := json.Marshal(got)
				if err != nil {
					t.Errorf("AppendChildren() marshal error = %v", err)
					return
				}
				b, err := json.Marshal(tt.want)
				if err != nil {
					t.Errorf("AppendChildren() marshal error = %v", err)
					return
				}

				if !(string(a) == string(b)) {
					t.Errorf("AppendChildren() got = %v, want %v", got, tt.want)
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
					BasicBlock: notionapi.BasicBlock{
						Object:         notionapi.ObjectTypeBlock,
						ID:             "some_id",
						Type:           notionapi.BlockTypeChildPage,
						CreatedTime:    &timestamp,
						LastEditedTime: &timestamp,
						CreatedBy: &notionapi.User{
							Object: "user",
							ID:     "some_id",
						},
						LastEditedBy: &notionapi.User{
							Object: "user",
							ID:     "some_id",
						},
						HasChildren: true,
					},
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
						RichText: []notionapi.RichText{
							{
								Text: &notionapi.Text{Content: "Hello"},
							},
						},
						Color: notionapi.ColorYellow.String(),
					},
				},
				want: &notionapi.ParagraphBlock{
					BasicBlock: notionapi.BasicBlock{
						Object:         notionapi.ObjectTypeBlock,
						ID:             "some_id",
						Type:           notionapi.BlockTypeParagraph,
						CreatedTime:    &timestamp,
						LastEditedTime: &timestamp,
					},
					Paragraph: notionapi.Paragraph{
						RichText: []notionapi.RichText{
							{
								Type: notionapi.ObjectTypeText,
								Text: &notionapi.Text{
									Content: "Hello",
								},
								Annotations: &notionapi.Annotations{Color: notionapi.ColorDefault},
								PlainText:   "Hello",
							},
						},
						Color: notionapi.ColorYellow.String(),
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

func TestBlockArrayUnmarshal(t *testing.T) {
	timestamp, err := time.Parse(time.RFC3339, "2021-11-04T02:09:00Z")
	if err != nil {
		t.Fatal(err)
	}

	var emoji notionapi.Emoji = "📌"
	var user *notionapi.User = &notionapi.User{
		Object: "user",
		ID:     "some_id",
	}
	t.Run("BlockArray", func(t *testing.T) {
		tests := []struct {
			name     string
			filePath string
			want     notionapi.Blocks
			wantErr  bool
			err      error
		}{
			{
				name:     "unmarshal",
				filePath: "testdata/block_array_unmarshal.json",
				want: notionapi.Blocks{
					&notionapi.CalloutBlock{
						BasicBlock: notionapi.BasicBlock{
							Object:         "block",
							ID:             "block1",
							Type:           "callout",
							CreatedTime:    &timestamp,
							LastEditedTime: &timestamp,
							CreatedBy:      user,
							LastEditedBy:   user,
						},
						Callout: notionapi.Callout{
							RichText: []notionapi.RichText{
								{
									Type: "text",
									Text: &notionapi.Text{
										Content: "This page is designed to be shared with students on the web. Click ",
									},
									Annotations: &notionapi.Annotations{
										Color: "default",
									},
									PlainText: "This page is designed to be shared with students on the web. Click ",
								}, {
									Type: "text",
									Text: &notionapi.Text{
										Content: "Share",
									},
									Annotations: &notionapi.Annotations{
										Code:  true,
										Color: "default",
									},
									PlainText: "Share",
								},
							},
							Icon: &notionapi.Icon{
								Type:  "emoji",
								Emoji: &emoji,
							},
							Color: notionapi.ColorBlue.String(),
						},
					},
					&notionapi.Heading1Block{
						BasicBlock: notionapi.BasicBlock{
							Object:         "block",
							ID:             "block2",
							Type:           "heading_1",
							CreatedTime:    &timestamp,
							LastEditedTime: &timestamp,
							CreatedBy:      user,
							LastEditedBy:   user,
						},
						Heading1: notionapi.Heading{
							RichText: []notionapi.RichText{
								{
									Type: "text",
									Text: &notionapi.Text{
										Content: "History 340",
									},
									Annotations: &notionapi.Annotations{
										Color: "default",
									},
									PlainText: "History 340",
								},
							},
							Color: notionapi.ColorBrownBackground.String(),
						},
					},
					&notionapi.ChildDatabaseBlock{
						BasicBlock: notionapi.BasicBlock{
							Object:         "block",
							ID:             "block3",
							Type:           "child_database",
							CreatedTime:    &timestamp,
							LastEditedTime: &timestamp,
							CreatedBy:      user,
							LastEditedBy:   user,
						},
						ChildDatabase: struct {
							Title string "json:\"title\""
						}{
							Title: "Required Texts",
						},
					},
					&notionapi.ColumnListBlock{
						BasicBlock: notionapi.BasicBlock{
							Object:         "block",
							ID:             "block4",
							Type:           "column_list",
							CreatedTime:    &timestamp,
							LastEditedTime: &timestamp,
							CreatedBy:      user,
							LastEditedBy:   user,
							HasChildren:    true,
						},
					},
					&notionapi.Heading3Block{
						BasicBlock: notionapi.BasicBlock{
							Object:         "block",
							ID:             "block5",
							Type:           "heading_3",
							CreatedTime:    &timestamp,
							LastEditedTime: &timestamp,
							CreatedBy:      user,
							LastEditedBy:   user,
						},
						Heading3: notionapi.Heading{
							RichText: []notionapi.RichText{
								{
									Type: "text",
									Text: &notionapi.Text{
										Content: "Assignment Submission",
									},
									Annotations: &notionapi.Annotations{
										Bold:  true,
										Color: "default",
									},
									PlainText: "Assignment Submission",
								},
							},
							Color: notionapi.ColorDefault.String(),
						},
					},
					&notionapi.ParagraphBlock{
						BasicBlock: notionapi.BasicBlock{
							Object:         "block",
							ID:             "block6",
							Type:           "paragraph",
							CreatedTime:    &timestamp,
							LastEditedTime: &timestamp,
							CreatedBy:      user,
							LastEditedBy:   user,
						},
						Paragraph: notionapi.Paragraph{
							RichText: []notionapi.RichText{
								{
									Type: "text",
									Text: &notionapi.Text{
										Content: "All essays and papers are due in lecture (due dates are listed on the schedule). No electronic copies will be accepted!",
									},
									Annotations: &notionapi.Annotations{
										Color: "default",
									},
									PlainText: "All essays and papers are due in lecture (due dates are listed on the schedule). No electronic copies will be accepted!",
								},
							},
							Color: notionapi.ColorRed.String(),
						},
					},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				data, err := ioutil.ReadFile(tt.filePath)
				if err != nil {
					t.Fatal(err)
				}
				got := make(notionapi.Blocks, 0)
				err = json.Unmarshal(data, &got)
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

func TestBlockUpdateRequest_MarshallJSON(t *testing.T) {
	tests := []struct {
		name    string
		req     *notionapi.BlockUpdateRequest
		want    []byte
		wantErr bool
	}{
		{
			name: "update todo checkbox",
			req: &notionapi.BlockUpdateRequest{
				ToDo: &notionapi.ToDo{Checked: false, RichText: make([]notionapi.RichText, 0)},
			},
			want: []byte(`{"to_do":{"rich_text":[],"checked":false}}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.req)
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
