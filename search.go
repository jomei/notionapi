package notionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SearchService interface {
	Do(context.Context, *SearchRequest) (*SearchResponse, error)
}

type SearchClient struct {
	apiClient *Client
}

// Do search https://developers.notion.com/reference/post-search
func (sc *SearchClient) Do(ctx context.Context, request *SearchRequest) (*SearchResponse, error) {
	res, err := sc.apiClient.request(ctx, http.MethodPost, "search", nil, request)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	var response SearchResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type SearchRequest struct {
	Query       string       `json:"query,omitempty"`
	Sort        *SortObject  `json:"sort,omitempty"`
	Filter      SearchFilter `json:"filter,omitempty"`
	StartCursor Cursor       `json:"start_cursor,omitempty"`
	PageSize    int          `json:"page_size,omitempty"`
}

type SearchResponse struct {
	Object     ObjectType `json:"object"`
	Results    []Object   `json:"results"`
	HasMore    bool       `json:"has_more"`
	NextCursor Cursor     `json:"next_cursor"`
}

func (sr *SearchResponse) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Object     ObjectType    `json:"object"`
		Results    []interface{} `json:"results"`
		HasMore    bool          `json:"has_more"`
		NextCursor Cursor        `json:"next_cursor"`
	}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	objects := make([]Object, len(tmp.Results))
	for i, rawObject := range tmp.Results {
		var o Object
		switch rawObject.(map[string]interface{})["object"].(string) {
		case ObjectTypeDatabase.String():
			o = &Database{}
		case ObjectTypePage.String():
			o = &Page{}
		default:
			return fmt.Errorf("unsupported object type %s", rawObject.(map[string]interface{})["object"].(string))
		}
		j, err := json.Marshal(rawObject)
		if err != nil {
			return err
		}

		err = json.Unmarshal(j, o)
		if err != nil {
			return err
		}
		objects[i] = o
	}

	*sr = SearchResponse{
		Object:     tmp.Object,
		Results:    objects,
		HasMore:    tmp.HasMore,
		NextCursor: tmp.NextCursor,
	}

	return nil
}
