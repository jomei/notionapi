package notionapi

import (
	"context"
	"encoding/json"
	"net/http"
)

type SearchService interface {
	Do(context.Context, SearchRequest) (*SearchResponse, error)
}

type SearchClient struct {
	apiClient *Client
}

// Do search https://developers.notion.com/reference/post-search
func (sc *SearchClient) Do(ctx context.Context, request SearchRequest) (*SearchResponse, error) {
	res, err := sc.apiClient.request(ctx, http.MethodPost, "search", nil, request)
	if err != nil {
		return nil, err
	}

	var response SearchResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type SearchRequest struct {
	Query       string        `json:"query,omitempty"`
	Sort        *SortObject   `json:"sort,omitempty"`
	Filter      *FilterObject `json:"filter,omitempty"`
	StartCursor Cursor        `json:"start_cursor,omitempty"`
	PageSize    int           `json:"page_size"`
}

type SearchResponse struct {
	Object     ObjectType    `json:"object"`
	Result     []BasicObject `json:"result"`
	HasMore    bool          `json:"has_more"`
	NextCursor Cursor        `json:"next_cursor"`
}
