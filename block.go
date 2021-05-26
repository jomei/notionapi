package notionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type BlockID string

func (bID BlockID) String() string {
	return string(bID)
}

type BlockType string

const (
	BlockTypeParagraph BlockType = "paragraph"
	BlockTypeHeading1  BlockType = "heading_1"
	BlockTypeHeading2  BlockType = "heading_2"
	BlockTypeHeading3  BlockType = "heading_3"

	BlockTypeBulletedListItem BlockType = "bulleted_list_item"
	BlockTypeNumberedListItem BlockType = "numbered_list_item"

	BlockTypeToDo        BlockType = "to_do"
	BlockTypeToggle      BlockType = "toggle"
	BlockTypeChildPage   BlockType = "child_page"
	BlockTypeUnsupported BlockType = "unsupported"
)

type BlockObject struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id"`
	CreatedTime    time.Time  `json:"created_time"`
	LastEditedTime time.Time  `json:"last_edited_time"`
	HasChildren    bool       `json:"has_children"`
}

func (c *Client) RetrieveBlockChildren(ctx context.Context, id BlockID, startCursor Cursor, pageSize int) ([]Object, error) {
	req, err := c.makeRetrieveBlockChildrenRequest(id, startCursor, pageSize)
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

	var response []Object
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil

}

func (c *Client) makeRetrieveBlockChildrenRequest(id BlockID, startCursor Cursor, pageSize int) (*http.Request, error) {
	reqURL := fmt.Sprintf("%s/%s/blocks/%s", ApiURL, ApiVersion, id.String())
	urlObj, err := url.Parse(reqURL)
	if err != nil {
		return nil, err
	}

	query := urlObj.Query()
	query.Add("start_cursor", startCursor.String())
	query.Add("page_size", strconv.Itoa(pageSize)) //TODO: empty values?

	urlObj.RawQuery = query.Encode()
	req, err := http.NewRequest(http.MethodGet, urlObj.String(), nil)
	if err != nil {
		return nil, err
	}

	req = c.addRequestHeaders(req)

	return req, nil
}
