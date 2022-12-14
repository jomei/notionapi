package notionapi_test

import (
	"context"
	"github.com/conduitio-labs/notionapi"
	"net/http"
	"reflect"
	"testing"
)

func TestUserClient(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			id         notionapi.UserID
			want       *notionapi.User
			wantErr    bool
			err        error
		}{
			{
				name:       "returns user by id",
				id:         "some_id",
				filePath:   "testdata/user_get.json",
				statusCode: http.StatusOK,
				want: &notionapi.User{
					Object:    notionapi.ObjectTypeUser,
					ID:        "some_id",
					Type:      notionapi.UserTypePerson,
					Name:      "John Doe",
					AvatarURL: "some.url",
					Person:    &notionapi.Person{Email: "some@email.com"},
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))

				got, err := client.User.Get(context.Background(), tt.id)
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

	t.Run("List", func(t *testing.T) {
		tests := []struct {
			name       string
			filePath   string
			statusCode int
			want       *notionapi.UsersListResponse
			wantErr    bool
			err        error
		}{
			{
				name:       "returns list of users",
				filePath:   "testdata/user_list.json",
				statusCode: http.StatusOK,
				want: &notionapi.UsersListResponse{
					Object: notionapi.ObjectTypeList,
					Results: []notionapi.User{
						{
							Object:    notionapi.ObjectTypeUser,
							ID:        "some_id",
							Type:      notionapi.UserTypePerson,
							Name:      "John Doe",
							AvatarURL: "some.url",
							Person:    &notionapi.Person{Email: "some@email.com"},
						},
						{
							Object: notionapi.ObjectTypeUser,
							ID:     "some_id",
							Type:   notionapi.UserTypeBot,
							Name:   "Test",
							Bot:    &notionapi.Bot{},
						},
					},
					HasMore: false,
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := newMockedClient(t, tt.filePath, tt.statusCode)
				client := notionapi.NewClient("some_token", notionapi.WithHTTPClient(c))
				got, err := client.User.List(context.Background(), nil)
				if (err != nil) != tt.wantErr {
					t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("List() got = %v, want %v", got, tt.want)
				}
			})
		}
	})
}
