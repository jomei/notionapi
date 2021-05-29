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
	Get(context.Context, PageID) (*Page, error)
	Create(context.Context, *PageCreateRequest) (*Page, error)
	Update(context.Context, PageID, map[string]BasicObject) (*Page, error)
}

type PageClient struct {
	apiClient *Client
}

// Get https://developers.notion.com/reference/get-page
func (pc *PageClient) Get(ctx context.Context, id PageID) (*Page, error) {
	res, err := pc.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("pages/%s", id.String()), nil, nil)
	if err != nil {
		return nil, err
	}

	return handlePageResponse(res)
}

// Create https://developers.notion.com/reference/post-page
func (pc *PageClient) Create(ctx context.Context, requestBody *PageCreateRequest) (*Page, error) {
	res, err := pc.apiClient.request(ctx, http.MethodPost, "pages", nil, requestBody)
	if err != nil {
		return nil, err
	}

	return handlePageResponse(res)
}

// Update https://developers.notion.com/reference/patch-page
func (pc *PageClient) Update(ctx context.Context, id PageID, properties map[string]BasicObject) (*Page, error) {
	res, err := pc.apiClient.request(ctx, http.MethodPatch, fmt.Sprintf("pages/%s", id.String()), nil, properties)
	if err != nil {
		return nil, err
	}

	return handlePageResponse(res)
}

type Page struct {
	Object         ObjectType `json:"object"`
	ID             ObjectID   `json:"id"`
	CreatedTime    time.Time  `json:"created_time"`
	LastEditedTime time.Time  `json:"last_edited_time"`
	Archived       bool       `json:"archived"`
	Properties     Properties `json:"properties"`
	Parent         Parent     `json:"parent"`
}

type ParentType string

type Parent struct {
	Type       ParentType `json:"type"`
	PageID     PageID     `json:"page_id,omitempty"`
	DatabaseID DatabaseID `json:"database_id,omitempty"`
}

type PageCreateRequest struct {
	Parent     Parent        `json:"parent"`
	Properties Properties    `json:"properties"`
	Children   []BlockObject `json:"children,omitempty"`
}

func handlePageResponse(res *http.Response) (*Page, error) {
	var response Page
	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
