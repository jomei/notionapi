package notionapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

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

func (c *Client) Search(ctx context.Context, request SearchRequest) (*SearchResponse, error) {
	req, err := c.makeSearchRequest(&request)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("http status: %d", res.StatusCode)
	}

	var response SearchResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) makeSearchRequest(request *SearchRequest) (*http.Request, error) {
	reqURL := fmt.Sprintf("%s/%s/search", ApiURL, ApiVersion)
	body, err := json.Marshal(request)
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
