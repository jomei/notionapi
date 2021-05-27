package notionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type BlockID string

func (bID BlockID) String() string {
	return string(bID)
}

type BlockObject struct {
	BasicObject
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id"`
	CreatedTime    time.Time  `json:"created_time"`
	LastEditedTime time.Time  `json:"last_edited_time"`
	HasChildren    bool       `json:"has_children"`
}

func (c *Client) BlockChildrenRetrieve(ctx context.Context, id BlockID, startCursor Cursor, pageSize int) ([]BasicObject, error) {
	queryParams := map[string]string{"start_cursor": startCursor.String(), "page_size": strconv.Itoa(pageSize)}
	res, err := c.request(ctx, http.MethodGet, fmt.Sprintf("blocks/%s", id.String()), queryParams, nil)
	if err != nil {
		return nil, err
	}

	var response []BasicObject
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type AppendBlockChildrenRequest struct {
	Children []BasicObject `json:"children"`
}

func (c *Client) BlockChildrenAppend(ctx context.Context, id BlockID, requestBody AppendBlockChildrenRequest) (*BlockObject, error) {
	res, err := c.request(ctx, http.MethodPost, fmt.Sprintf("blocks/%s", id.String()), nil, requestBody)
	if err != nil {
		return nil, err
	}

	var response BlockObject
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
