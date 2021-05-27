package notionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PageID string

func (pID PageID) String() string {
	return string(pID)
}

type PageService interface {
	Retrieve(context.Context, PageID) (*PageObject, error)
	Create(context.Context, PageCreateRequest) (*PageObject, error)
	Update(context.Context, PageID, map[PropertyName]BasicObject) (*PageObject, error)
}

type PageClient struct {
	apiClient *Client
}

func (pc *PageClient) Retrieve(ctx context.Context, id PageID) (*PageObject, error) {
	res, err := pc.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("pages/%s", id.String()), nil, nil)
	if err != nil {
		return nil, err
	}

	return handlePageResponse(res)
}

func (pc *PageClient) Create(ctx context.Context, requestBody PageCreateRequest) (*PageObject, error) {
	res, err := pc.apiClient.request(ctx, http.MethodPost, "pages", nil, requestBody)
	if err != nil {
		return nil, err
	}

	return handlePageResponse(res)
}

func (pc *PageClient) Update(ctx context.Context, id PageID, properties map[PropertyName]BasicObject) (*PageObject, error) {
	res, err := pc.apiClient.request(ctx, http.MethodPatch, fmt.Sprintf("pages/%s", id.String()), nil, properties)
	if err != nil {
		return nil, err
	}

	return handlePageResponse(res)
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

type PageCreateRequest struct {
	Parent     Parent                       `json:"parent"`
	Properties map[PropertyName]BasicObject `json:"properties"`
	Children   []BlockObject                `json:"children"`
}

func handlePageResponse(res *http.Response) (*PageObject, error) {
	var response PageObject
	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
