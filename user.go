package notionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type UserID string

func (uID UserID) String() string {
	return string(uID)
}

type UserService interface {
	Get(context.Context, UserID) (*User, error)
	List(context.Context, Cursor, int) (*UsersListResponse, error)
}

type UserClient struct {
	apiClient *Client
}

// Get https://developers.notion.com/reference/get-user
func (uc *UserClient) Get(ctx context.Context, id UserID) (*User, error) {
	res, err := uc.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("users/%s", id.String()), nil, nil)
	if err != nil {
		return nil, err
	}

	var response User
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// List https://developers.notion.com/reference/get-users
func (uc *UserClient) List(ctx context.Context, startCursor Cursor, pageSize int) (*UsersListResponse, error) {
	queryParams := map[string]string{"start_cursor": startCursor.String(), "page_size": strconv.Itoa(pageSize)}
	res, err := uc.apiClient.request(ctx, http.MethodGet, "users", queryParams, nil)
	if err != nil {
		return nil, err
	}

	var response UsersListResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type UserType string

type User struct {
	Object    ObjectType `json:"object"`
	ID        UserID     `json:"id"`
	Type      UserType   `json:"type"`
	Name      string     `json:"name"`
	AvatarURL string     `json:"avatar_url"`
	Person    *Person    `json:"person"`
	Bot       *BotObject `json:"bot"`
}

type Person struct {
	Email string `json:"email"`
}

type BotObject struct{}

type UsersListResponse struct {
	Results    []User `json:"results"`
	HasMore    bool   `json:"has_more"`
	NextCursor Cursor `json:"next_cursor"`
}
