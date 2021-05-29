package notionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DatabaseID string

func (dID DatabaseID) String() string {
	return string(dID)
}

type DatabaseService interface {
	Get(context.Context, DatabaseID) (*Database, error)
	List(context.Context, *Pagination) (*DatabaseListResponse, error)
	Query(context.Context, DatabaseID, *DatabaseQueryRequest) (*DatabaseQueryResponse, error)
}

type DatabaseClient struct {
	apiClient *Client
}

// Get https://developers.notion.com/reference/get-database
func (dc *DatabaseClient) Get(ctx context.Context, id DatabaseID) (*Database, error) {
	res, err := dc.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("databases/%s", id.String()), nil, nil)
	if err != nil {
		return nil, err
	}
	var response Database

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// List https://developers.notion.com/reference/get-databases
func (dc *DatabaseClient) List(ctx context.Context, pagination *Pagination) (*DatabaseListResponse, error) {
	res, err := dc.apiClient.request(ctx, http.MethodGet, "databases", pagination.ToQuery(), nil)
	if err != nil {
		return nil, err
	}

	var response DatabaseListResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Query https://developers.notion.com/reference/post-database-query
func (dc *DatabaseClient) Query(ctx context.Context, id DatabaseID, requestBody *DatabaseQueryRequest) (*DatabaseQueryResponse, error) {
	res, err := dc.apiClient.request(ctx, http.MethodPost, fmt.Sprintf("databases/%s/query", id.String()), nil, requestBody)
	if err != nil {
		return nil, err
	}

	var response DatabaseQueryResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type Database struct {
	Object         ObjectType `json:"object"`
	ID             ObjectID   `json:"id"`
	CreatedTime    time.Time  `json:"created_time"`
	LastEditedTime time.Time  `json:"last_edited_time"`
	Title          Paragraph  `json:"title"`
	Properties     Properties `json:"properties"`
}

type DatabaseListResponse struct {
	Object     ObjectType `json:"object"`
	Results    []Database `json:"results"`
	NextCursor string     `json:"next_cursor"`
	HasMore    bool       `json:"has_more"`
}

type DatabaseQueryRequest struct {
	Filter      Filter       `json:"filter,omitempty"`
	Sorts       []SortObject `json:"sorts"`
	StartCursor Cursor       `json:"start_cursor,omitempty"`
	PageSize    int          `json:"page_size,omitempty"`
}

type DatabaseQueryResponse struct {
	Object     ObjectType `json:"object"`
	Results    []Page     `json:"results"`
	HasMore    bool       `json:"has_more"`
	NextCursor Cursor     `json:"next_cursor"`
}
