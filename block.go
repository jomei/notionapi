package notionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BlockID string

func (bID BlockID) String() string {
	return string(bID)
}

type BlockService interface {
	GetChildren(context.Context, BlockID, *Pagination) ([]BasicObject, error)
	AppendChildren(context.Context, BlockID, AppendBlockChildrenRequest) (*BlockObject, error)
}

type BlockClient struct {
	apiClient *Client
}

// GetChildren https://developers.notion.com/reference/get-block-children
func (bc *BlockClient) GetChildren(ctx context.Context, id BlockID, pagination *Pagination) ([]BasicObject, error) {
	res, err := bc.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("blocks/%s", id.String()), pagination.ToQuery(), nil)
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

// AppendChildren https://developers.notion.com/reference/patch-block-children
func (bc *BlockClient) AppendChildren(ctx context.Context, id BlockID, requestBody AppendBlockChildrenRequest) (*BlockObject, error) {
	res, err := bc.apiClient.request(ctx, http.MethodPost, fmt.Sprintf("blocks/%s", id.String()), nil, requestBody)
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

type BlockObject struct {
	BasicObject
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id"`
	CreatedTime    time.Time  `json:"created_time"`
	LastEditedTime time.Time  `json:"last_edited_time"`
	HasChildren    bool       `json:"has_children"`
}

type AppendBlockChildrenRequest struct {
	Children []BasicObject `json:"children"`
}
