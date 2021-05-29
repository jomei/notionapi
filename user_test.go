package notionapi_test

import (
	"context"
	"github.com/jomei/notionapi"
	"reflect"
	"testing"
)

func TestUserClient(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		tests := []struct {
			name     string
			filePath string
			id       notionapi.UserID
			want     *notionapi.User
			wantErr  bool
			err      error
		}{
			{
				name:     "returns user by id",
				id:       "some_id",
				filePath: "testdata/user_get.json",
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
				c := newMockedClient(t, tt.filePath)
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
}
