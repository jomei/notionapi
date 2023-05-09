package notionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type DatabaseID string

func (dID DatabaseID) String() string {
	return string(dID)
}

type DatabaseService interface {
	Get(context.Context, DatabaseID) (*Database, error)
	Query(context.Context, DatabaseID, *DatabaseQueryRequest) (*DatabaseQueryResponse, error)
	Update(context.Context, DatabaseID, *DatabaseUpdateRequest) (*Database, error)
	Create(ctx context.Context, request *DatabaseCreateRequest) (*Database, error)
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

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	var response Database

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

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	var response DatabaseQueryResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Update https://developers.notion.com/reference/update-a-database
func (dc *DatabaseClient) Update(ctx context.Context, id DatabaseID, requestBody *DatabaseUpdateRequest) (*Database, error) {
	res, err := dc.apiClient.request(ctx, http.MethodPatch, fmt.Sprintf("databases/%s", id.String()), nil, requestBody)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	var response Database
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

// Create https://developers.notion.com/reference/create-a-database
func (dc *DatabaseClient) Create(ctx context.Context, requestBody *DatabaseCreateRequest) (*Database, error) {
	res, err := dc.apiClient.request(ctx, http.MethodPost, "databases", nil, requestBody)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	var response Database
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
	CreatedBy      User       `json:"created_by,omitempty"`
	LastEditedBy   User       `json:"last_edited_by,omitempty"`
	Title          []RichText `json:"title"`
	Parent         Parent     `json:"parent"`
	URL            string     `json:"url"`
	// Properties is a map of property configurations that defines what Page.Properties each page of the database can use
	Properties  PropertyConfigs `json:"properties"`
	Description []RichText      `json:"description"`
	IsInline    bool            `json:"is_inline"`
	Archived    bool            `json:"archived"`
	Icon        *Icon           `json:"icon,omitempty"`
	Cover       *Image          `json:"cover,omitempty"`
}

func (db *Database) GetObject() ObjectType {
	return db.Object
}

type DatabaseListResponse struct {
	Object     ObjectType `json:"object"`
	Results    []Database `json:"results"`
	NextCursor string     `json:"next_cursor"`
	HasMore    bool       `json:"has_more"`
}

type DatabaseQueryRequest struct {
	Filter      Filter
	Sorts       []SortObject `json:"sorts,omitempty"`
	StartCursor Cursor       `json:"start_cursor,omitempty"`
	PageSize    int          `json:"page_size,omitempty"`
}

func (qr *DatabaseQueryRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sorts       []SortObject `json:"sorts,omitempty"`
		StartCursor Cursor       `json:"start_cursor,omitempty"`
		PageSize    int          `json:"page_size,omitempty"`
		Filter      interface{}  `json:"filter,omitempty"`
	}{
		Sorts:       qr.Sorts,
		StartCursor: qr.StartCursor,
		PageSize:    qr.PageSize,
		Filter:      qr.Filter,
	})
}

type DatabaseQueryResponse struct {
	Object     ObjectType `json:"object"`
	Results    []Page     `json:"results"`
	HasMore    bool       `json:"has_more"`
	NextCursor Cursor     `json:"next_cursor"`
}

type DatabaseUpdateRequest struct {
	Title      []RichText      `json:"title,omitempty"`
	Properties PropertyConfigs `json:"properties,omitempty"`
}

type DatabaseCreateRequest struct {
	Parent     Parent          `json:"parent"`
	Title      []RichText      `json:"title"`
	Properties PropertyConfigs `json:"properties"`
	IsInline   bool            `json:"is_inline"`
}
