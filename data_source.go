package notionapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type DataSourceID string

func (dsID DataSourceID) String() string {
	return string(dsID)
}

type DataSourceService interface {
	Create(ctx context.Context, request *DataSourceCreateRequest) (*DataSource, error)
	Query(ctx context.Context, id DataSourceID, request *DataSourceQueryRequest) (*DataSourceQueryResponse, error)
	Get(ctx context.Context, id DataSourceID) (*DataSource, error)
	Update(ctx context.Context, id DataSourceID, request *DataSourceUpdateRequest) (*DataSource, error)
}

type DataSourceClient struct {
	apiClient *Client
}

// Create a data source.
// https://developers.notion.com/reference/create-a-data-source
func (c *DataSourceClient) Create(ctx context.Context, requestBody *DataSourceCreateRequest) (*DataSource, error) {
	res, err := c.apiClient.request(ctx, http.MethodPost, "data_sources", nil, requestBody)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	var response DataSource
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

type DataSourceCreateRequest struct {
	Parent     Parent          `json:"parent"`
	Properties PropertyConfigs `json:"properties"`
	Title      []RichText      `json:"title"`
	IsInline   bool            `json:"is_inline,omitempty"`
}

// Query gets a list of Pages contained in the data source, filtered and ordered
// according to the filter conditions and sort criteria provided in the request.
// The response may contain fewer than page_size of results. If the response
// includes a next_cursor value, refer to the pagination reference for details
// about how to use a cursor to iterate through the list.
//
// Filters are similar to the filters provided in the Notion UI where the set of
// filters and filter groups chained by "And" in the UI is equivalent to having
// each filter in the array of the compound "and" filter. A similar set of
// filters chained by "Or" in the UI would be represented as filters in the
// array of the "or" compound filter.
//
// Filters operate on data source properties and can be combined. If no filter is
// provided, all the pages in the data source will be returned with pagination.
//
// See https://developers.notion.com/reference/post-data-source-query.
func (c *DataSourceClient) Query(ctx context.Context, id DataSourceID, requestBody *DataSourceQueryRequest) (*DataSourceQueryResponse, error) {
	res, err := c.apiClient.request(ctx, http.MethodPost, fmt.Sprintf("data_sources/%s/query", id.String()), nil, requestBody)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	var response DataSourceQueryResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Get a data source.
// https://developers.notion.com/reference/get-data-source
func (c *DataSourceClient) Get(ctx context.Context, id DataSourceID) (*DataSource, error) {
	if id == "" {
		return nil, errors.New("empty data source id")
	}

	res, err := c.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("data_sources/%s", id.String()), nil, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	var response DataSource
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Update a data source.
// https://developers.notion.com/reference/update-a-data-source
func (c *DataSourceClient) Update(ctx context.Context, id DataSourceID, requestBody *DataSourceUpdateRequest) (*DataSource, error) {
	res, err := c.apiClient.request(ctx, http.MethodPatch, fmt.Sprintf("data_sources/%s", id.String()), nil, requestBody)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	var response DataSource
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

type DataSourceUpdateRequest struct {
	Title      []RichText      `json:"title,omitempty"`
	Properties PropertyConfigs `json:"properties,omitempty"`
}

// DataSource object.
type DataSource struct {
	Object         ObjectType      `json:"object"`
	ID             ObjectID        `json:"id"`
	CreatedTime    time.Time       `json:"created_time"`
	LastEditedTime time.Time       `json:"last_edited_time"`
	CreatedBy      User            `json:"created_by,omitempty"`
	LastEditedBy   User            `json:"last_edited_by,omitempty"`
	Title          []RichText      `json:"title"`
	Parent         Parent          `json:"parent"`
	DatabaseParent *Parent         `json:"database_parent,omitempty"`
	URL            string          `json:"url"`
	PublicURL      string          `json:"public_url"`
	Properties     PropertyConfigs `json:"properties"`
	IsInline       bool            `json:"is_inline"`
	Archived       bool            `json:"archived"`
	Icon           *Icon           `json:"icon,omitempty"`
	Cover          *Image          `json:"cover,omitempty"`
	Description    []RichText      `json:"description,omitempty"`
}

func (ds *DataSource) GetObject() ObjectType {
	return ds.Object
}

type DataSourceQueryRequest struct {
	Filter      Filter       `json:"filter,omitempty"`
	Sorts       []SortObject `json:"sorts,omitempty"`
	StartCursor Cursor       `json:"start_cursor,omitempty"`
	PageSize    int          `json:"page_size,omitempty"`
}

func (qr *DataSourceQueryRequest) MarshalJSON() ([]byte, error) {
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

type DataSourceQueryResponse struct {
	Object     ObjectType `json:"object"`
	Results    []Page     `json:"results"`
	HasMore    bool       `json:"has_more"`
	NextCursor Cursor     `json:"next_cursor"`
}
