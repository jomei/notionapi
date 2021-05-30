package notionapi_test

import (
	"context"
	"github.com/jomei/notionapi"
	"testing"
)

func TestSearchClient(t *testing.T) {
	t.Run("Do", func(t *testing.T) {
		tests := []struct {
			name     string
			filePath string
			request  *notionapi.SearchRequest
			wantErr  bool
			err      error
		}{
			{
				name:     "returns search result",
				filePath: "testdata/search.json",
				request: &notionapi.SearchRequest{
					Query: "Hel",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))

				got, err := client.Search.Do(context.Background(), tt.request)

				if (err != nil) != tt.wantErr {
					t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if got == nil {
					t.Errorf("Search result is nil")
				}
			})
		}
	})
}
