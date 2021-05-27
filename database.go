package notionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type DatabaseID string

func (dID DatabaseID) String() string {
	return string(dID)
}

func (c *Client) DBRetrieve(ctx context.Context, id DatabaseID) (*DatabaseObject, error) {
	res, err := c.request(ctx, http.MethodGet, fmt.Sprintf("databases/%s", id.String()), nil, nil)
	if err != nil {
		return nil, err
	}
	var response DatabaseObject
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type PropertyName string

func (pn PropertyName) String() string {
	return string(pn)
}

type DatabaseObject struct {
	Object         ObjectType                   `json:"object"`
	ID             ObjectID                     `json:"id"`
	CreatedTime    time.Time                    `json:"created_time"` //TODO: format
	LastEditedTime time.Time                    `json:"last_edited_time"`
	Title          TextObject                   `json:"title"`
	Properties     map[PropertyName]BasicObject `json:"properties"`
}

func (c *Client) DBList(ctx context.Context, startCursor Cursor, pageSize int) (*DBListResponse, error) {
	queryParams := map[string]string{"start_cursor": startCursor.String(), "page_size": strconv.Itoa(pageSize)}
	res, err := c.request(ctx, http.MethodGet, "/databases", queryParams, nil)
	if err != nil {
		return nil, err
	}

	var response DBListResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil

}

type DBListResponse struct {
	Results    []DatabaseObject `json:"results"`
	NextCursor string           `json:"next_cursor"`
	HasMore    bool             `json:"has_more"`
}

type DBQueryRequest struct {
	Filter      *FilterObject `json:"filter,omitempty"`
	Sorts       []SortObject  `json:"sorts"`
	StartCursor Cursor        `json:"start_cursor,omiempty"`
	PageSize    int           `json:"page_size"`
}

func (c *Client) DBQuery(ctx context.Context, id DatabaseID, requestBody DBQueryRequest) (*DBQueryResponse, error) {
	res, err := c.request(ctx, http.MethodPost, fmt.Sprintf("databases/%s", id.String()), nil, requestBody)
	if err != nil {
		return nil, err
	}

	var response DBQueryResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type DBQueryResponse struct {
	Object     ObjectType   `json:"object"`
	Results    []PageObject `json:"results"`
	HasMore    bool         `json:"has_more"`
	NextCursor Cursor       `json:"next_cursor"`
}
