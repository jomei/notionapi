package notionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type UserID string

func (uID UserID) String() string {
	return string(uID)
}

type UserService interface {
	Get(context.Context, UserID) (*User, error)
	List(context.Context, *Pagination) (*UsersListResponse, error)
	Me(context.Context) (*User, error)
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

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	var response User
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// List https://developers.notion.com/reference/get-users
func (uc *UserClient) List(ctx context.Context, pagination *Pagination) (*UsersListResponse, error) {
	res, err := uc.apiClient.request(ctx, http.MethodGet, "users", pagination.ToQuery(), nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	var response UsersListResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Me https://developers.notion.com/reference/get-self
func (uc *UserClient) Me(ctx context.Context) (*User, error) {
	res, err := uc.apiClient.request(ctx, http.MethodGet, "users/me", nil, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	var response User
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type UserType string

type User struct {
	Object    ObjectType `json:"object,omitempty"`
	ID        UserID     `json:"id"`
	Type      UserType   `json:"type,omitempty"`
	Name      string     `json:"name,omitempty"`
	AvatarURL string     `json:"avatar_url,omitempty"`
	Person    *Person    `json:"person,omitempty"`
	Bot       *Bot       `json:"bot,omitempty"`
}

type Person struct {
	Email string `json:"email"`
}

type Bot struct {
	Owner         Owner  `json:"owner"`
	WorkspaceName string `json:"workspace_name"`
}

type Owner struct {
	Type      string `json:"type"`
	Workspace bool   `json:"workspace"`
}

type UsersListResponse struct {
	Object     ObjectType `json:"object"`
	Results    []User     `json:"results"`
	HasMore    bool       `json:"has_more"`
	NextCursor Cursor     `json:"next_cursor"`
}
