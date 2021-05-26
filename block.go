package notionapi

import (
	"bytes"
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

type BlockObject struct {
	BasicObject
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id"`
	CreatedTime    time.Time  `json:"created_time"`
	LastEditedTime time.Time  `json:"last_edited_time"`
	HasChildren    bool       `json:"has_children"`
}

func (c *Client) BlockChildrenRetrieve(ctx context.Context, id BlockID, startCursor Cursor, pageSize int) ([]BasicObject, error) {
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

	var response []BasicObject
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

type AppendBlockChildrenRequest struct {
	Children []BasicObject `json:"children"`
}

func (c *Client) BlockChildrenAppend(ctx context.Context, id BlockID, requestBody AppendBlockChildrenRequest) (*BlockObject, error) {
	req, err := c.makeAppendBlockChildrenRequest(id, &requestBody)
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

	var response BlockObject
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) makeAppendBlockChildrenRequest(id BlockID, requestBody *AppendBlockChildrenRequest) (*http.Request, error) {
	reqURL := fmt.Sprintf("%s/%s/blocks/%s/children", ApiURL, ApiVersion, id.String())
	body, err := json.Marshal(requestBody)
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
