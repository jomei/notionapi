package notionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
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
	req, err := c.makeUserRetrieveRequest(id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("http status: %d", res.StatusCode)
	}

	var response UserObject
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) makeUserRetrieveRequest(id UserID) (*http.Request, error) {
	reqURL := fmt.Sprintf("%s/%s/users/%s", ApiURL, ApiVersion, id.String())
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req = c.addRequestHeaders(req)

	return req, nil
}

type UsersListResponse struct {
	Results []UserObject `json:"results"`
	HasMore bool         `json:"has_more"`
}

func (c *Client) UsersList(ctx context.Context, startCursor Cursor, pageSize int) (*UsersListResponse, error) {
	req, err := c.makeUsersListRequest(startCursor, pageSize)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("http status: %d", res.StatusCode)
	}

	var response UsersListResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) makeUsersListRequest(startCursor Cursor, pageSize int) (*http.Request, error) {
	reqURL := fmt.Sprintf("%s/%s/users", ApiURL, ApiVersion)
	urlObj, err := url.Parse(reqURL)
	if err != nil {
		return nil, err
	}

	query := urlObj.Query()
	query.Add("start_cursor", startCursor.String())
	query.Add("page_size", strconv.Itoa(pageSize)) //TODO: empty values?

	urlObj.RawQuery = query.Encode()
	req, err := http.NewRequest(http.MethodGet, urlObj.String(), nil)
	if err != nil {
		return nil, err
	}

	req = c.addRequestHeaders(req)

	return req, nil
}
