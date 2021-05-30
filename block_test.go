package notionapi_test

import (
	"context"
	"github.com/jomei/notionapi"
	"testing"
)

func TestBlockClient(t *testing.T) {
	t.Run("GetChildren", func(t *testing.T) {
		tests := []struct {
			name     string
			filePath string
			id       notionapi.BlockID
			len      int
			wantErr  bool
			err      error
		}{
			{
				name:     "returns blocks by id of parent block",
				id:       "d1c2fdf4-9f12-46cc-b168-1ed1bcb732d8",
				filePath: "testdata/block_get_children.json",
				len:      2,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))
				got, err := client.Block.GetChildren(context.Background(), tt.id, nil)

				if (err != nil) != tt.wantErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if tt.len != len(got.Results) {
					t.Errorf("GetChildren got %d, want: %d", len(got.Results), tt.len)
				}
			})
		}
	})
}
