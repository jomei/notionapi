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

// Retrieves a Page object using the ID specified.
//
// Responses contains page properties, not page content. To fetch page content,
// use the Retrieve block children endpoint.
//
// Page properties are limited to up to 25 references per page property. To
// retrieve data related to properties that have more than 25 references, use
// the Retrieve a page property endpoint.
//
// See https://developers.notion.com/reference/get-page
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

// Creates a new page that is a child of an existing page or database.
//
// If the new page is a child of an existing page,title is the only valid
// property in the properties body param.
//
// If the new page is a child of an existing database, the keys of the
// properties object body param must match the parent database's properties.
//
// This endpoint can be used to create a new page with or without content using
// the children option. To add content to a page after creating it, use the
// Append block children endpoint.
//
// Returns a new page object.
//
// See https://developers.notion.com/reference/post-page
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

// Updates the properties of a page in a database. The properties body param of
// this endpoint can only be used to update the properties of a page that is a
// child of a database. The page’s properties schema must match the parent
// database’s properties.
//
// This endpoint can be used to update any page icon or cover, and can be used
// to archive or restore any page.
//
// To add page content instead of page properties, use the append block children
// endpoint. The page_id can be passed as the block_id when adding block
// children to the page.
//
// Returns the updated page object.
//
// See https://developers.notion.com/reference/patch-page
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

// The Page object contains the page property values of a single Notion page.
//
// See https://developers.notion.com/reference/page
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

// Pages, databases, and blocks are either located inside other pages,
// databases, and blocks, or are located at the top level of a workspace. This
// location is known as the "parent". Parent information is represented by a
// consistent parent object throughout the API.
//
// See https://developers.notion.com/reference/parent-object
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
