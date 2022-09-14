package notionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	Update(context.Context, PageID, *PageUpdateRequest) (*Page, error)
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

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	return handlePageResponse(res)
}

// Create https://developers.notion.com/reference/post-page
func (pc *PageClient) Create(ctx context.Context, requestBody *PageCreateRequest) (*Page, error) {
	res, err := pc.apiClient.request(ctx, http.MethodPost, "pages", nil, requestBody)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	return handlePageResponse(res)
}

type PageUpdateRequest struct {
	Properties Properties `json:"properties"`
	Archived   bool       `json:"archived"`
	Icon       *Icon      `json:"icon,omitempty"`
	Cover      *Image     `json:"cover,omitempty"`
}

// Update https://developers.notion.com/reference/patch-page
func (pc *PageClient) Update(ctx context.Context, id PageID, request *PageUpdateRequest) (*Page, error) {
	res, err := pc.apiClient.request(ctx, http.MethodPatch, fmt.Sprintf("pages/%s", id.String()), nil, request)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	return handlePageResponse(res)
}

type Page struct {
	Object         ObjectType `json:"object"`
	ID             ObjectID   `json:"id"`
	CreatedTime    time.Time  `json:"created_time"`
	LastEditedTime time.Time  `json:"last_edited_time"`
	CreatedBy      User       `json:"created_by,omitempty"`
	LastEditedBy   User       `json:"last_edited_by,omitempty"`
	Archived       bool       `json:"archived"`
	Properties     Properties `json:"properties"`
	Parent         Parent     `json:"parent"`
	URL            string     `json:"url"`
	Icon           *Icon      `json:"icon,omitempty"`
	Cover          *Image     `json:"cover,omitempty"`
}

func (p *Page) GetObject() ObjectType {
	return p.Object
}

type ParentType string

// Ref: https://developers.notion.com/reference/parent-object
type Parent struct {
	Type       ParentType `json:"type,omitempty"`
	PageID     PageID     `json:"page_id,omitempty"`
	DatabaseID DatabaseID `json:"database_id,omitempty"`
	BlockID    BlockID    `json:"block_id,omitempty"`
	Workspace  bool       `json:"workspace,omitempty"`
}

type PageCreateRequest struct {
	Parent     Parent     `json:"parent"`
	Properties Properties `json:"properties"`
	Children   []Block    `json:"children,omitempty"`
	Icon       *Icon      `json:"icon,omitempty"`
	Cover      *Image     `json:"cover,omitempty"`
}

func handlePageResponse(res *http.Response) (*Page, error) {
	var response Page
	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
