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

type UserObject struct {
	Object    ObjectType    `json:"object"`
	ID        UserID        `json:"id"`
	Type      ObjectType    `json:"type"`
	Name      string        `json:"name"`
	AvatarURL string        `json:"avatar_url"`
	Person    *PersonObject `json:"person"`
	Bot       *BotObject    `json:"bot"`
}

type PersonObject struct {
	Email string `json:"email"`
}

type BotObject struct{}

func (c *Client) UserRetrieve(ctx context.Context, id UserID) (*UserObject, error) {
	res, err := c.request(ctx, http.MethodGet, fmt.Sprintf("users/%s", id.String()), nil, nil)
	if err != nil {
		return nil, err
	}

	var response UserObject
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type UsersListResponse struct {
	Results []UserObject `json:"results"`
	HasMore bool         `json:"has_more"`
}

func (c *Client) UsersList(ctx context.Context, startCursor Cursor, pageSize int) (*UsersListResponse, error) {
	queryParams := map[string]string{"start_cursor": startCursor.String(), "page_size": strconv.Itoa(pageSize)}
	res, err := c.request(ctx, http.MethodGet, "users", queryParams, nil)
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
