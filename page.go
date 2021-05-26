package notionapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type PageID string

func (pID PageID) String() string {
	return string(pID)
}

type PageObject struct {
	Object         ObjectType                   `json:"object"`
	ID             ObjectID                     `json:"id"`
	CreatedTime    time.Time                    `json:"created_time"` // TODO: format
	LastEditedTime time.Time                    `json:"last_edited_time"`
	Archived       bool                         `json:"archived"`
	Properties     map[PropertyName]BasicObject `json:"properties"`
	Parent         Parent                       `json:"parent"`
}

type Parent struct {
	Type       ObjectType `json:"type"`
	PageID     PageID     `json:"page_id,omitempty"`
	DatabaseID DatabaseID `json:"database_id,omitempty"`
}

func (c *Client) PageRetrieve(ctx context.Context, id PageID) (*PageObject, error) {
	req, err := c.makePageRetrieveRequest(id)
	if err != nil {
		return nil, err
	}

	return c.doPageRequest(ctx, req)
}

func (c *Client) makePageRetrieveRequest(id PageID) (*http.Request, error) {
	reqURL := fmt.Sprintf("%s/%s/pages/%s", ApiURL, ApiVersion, id.String())
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req = c.addRequestHeaders(req)

	return req, nil
}

type PageCreateRequest struct {
	Parent     Parent                       `json:"parent"`
	Properties map[PropertyName]BasicObject `json:"properties"`
	Children   []BlockType                  `json:"children"`
}

func (c *Client) PageCreate(ctx context.Context, requestBody PageCreateRequest) (*PageObject, error) {
	req, err := c.makePageCreateRequest(requestBody)
	if err != nil {
		return nil, err
	}

	return c.doPageRequest(ctx, req)
}

func (c *Client) makePageCreateRequest(requestBody PageCreateRequest) (*http.Request, error) {
	reqURL := fmt.Sprintf("%s/%s/pages", ApiURL, ApiVersion)
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req = c.addRequestHeaders(req)

	return req, nil
}

func (c *Client) PageUpdate(ctx context.Context, id PageID, properties map[PropertyName]BasicObject) (*PageObject, error) {
	req, err := c.makePageUpdateRequest(id, properties)
	if err != nil {
		return nil, err
	}
	return c.doPageRequest(ctx, req)
}

func (c *Client) makePageUpdateRequest(id PageID, properties map[PropertyName]BasicObject) (*http.Request, error) {
	reqURL := fmt.Sprintf("%s/%s/pages/%s", ApiURL, ApiVersion, id.String())
	body, err := json.Marshal(properties)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPatch, reqURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req = c.addRequestHeaders(req)

	return req, nil
}

func (c *Client) doPageRequest(ctx context.Context, req *http.Request) (*PageObject, error) {
	req = req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("http status: %d", res.StatusCode)
	}

	var response PageObject
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
