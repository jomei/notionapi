package notionapi_test

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/jomei/notionapi"
)

func TestAuthenticationClient(t *testing.T) {
	t.Run("CreateToken", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			request    *notionapi.TokenCreateRequest
			want       *notionapi.TokenCreateResponse
			wantErr    bool
			err        error
		}{
			{
				name:       "Creates token",
				filePath:   "testdata/create_token.json",
				statusCode: http.StatusOK,
				request: &notionapi.TokenCreateRequest{
					Code:        "code1",
					GrantType:   "authorization_code",
					RedirectUri: "www.example.com",
				},
				want: &notionapi.TokenCreateResponse{
					AccessToken:          "token1",
					BotId:                "bot1",
					DuplicatedTemplateId: "template_id1",
					WorkspaceIcon:        "🎉",
					WorkspaceId:          "workspaceid_1",
					WorkspaceName:        "workspace_1",
				},
				wantErr: false,
				err:     nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))
				got, err := client.Authentication.CreateToken(context.Background(), tt.request)

				if (err != nil) != tt.wantErr {
					t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Query() got = %v, want %v", got, tt.want)
				}
			})
		}
	})
}
