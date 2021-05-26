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

type DatabaseID string

func (dID DatabaseID) String() string {
	return string(dID)
}

func (c *Client) DBRetrieve(ctx context.Context, id DatabaseID) (*DatabaseObject, error) {
	req, err := c.makeDBRetrieveRequest(id)
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

	var response DatabaseObject
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) makeDBRetrieveRequest(id DatabaseID) (*http.Request, error) {
	reqURL := fmt.Sprintf("%s/%s/databases/%s", ApiURL, ApiVersion, id.String())
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	req = c.addRequestHeaders(req)

	return req, nil
}

type PropertyName string

func (pn PropertyName) String() string {
	return string(pn)
}

type DatabaseObject struct {
	Object         ObjectType                   `json:"object"`
	ID             ObjectID                     `json:"id"`
	CreatedTime    time.Time                    `json:"created_time"` //TODO: format
	LastEditedTime time.Time                    `json:"last_edited_time"`
	Title          TextObject                   `json:"title"`
	Properties     map[PropertyName]BasicObject `json:"properties"`
}

func (c *Client) DBList(ctx context.Context, startCursor Cursor, pageSize int) (*DBListResponse, error) {
	req, err := c.makeDBListRequest(startCursor, pageSize)
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

	var response DBListResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil

}

func (c *Client) makeDBListRequest(startCursor Cursor, pageSize int) (*http.Request, error) {
	reqURL := fmt.Sprintf("%s/%s/databases", ApiURL, ApiVersion)
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

type DBListResponse struct {
	Results    []DatabaseObject `json:"results"`
	NextCursor string           `json:"next_cursor"`
	HasMore    bool             `json:"has_more"`
}

type DBQueryRequest struct {
	Filter      FilterObject `json:"filter"`
	Sorts       []SortObject `json:"sorts"`
	StartCursor Cursor       `json:"start_cursor"`
	PageSize    int          `json:"page_size"`
}

func (c *Client) DBQuery(ctx context.Context, id DatabaseID, requestBody DBQueryRequest) (*DBQueryResponse, error) {
	req, err := c.makeDBQueryRequest(id, requestBody)
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

	var response DBQueryResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) makeDBQueryRequest(id DatabaseID, requestBody DBQueryRequest) (*http.Request, error) {
	reqURL := fmt.Sprintf("%s/%s/databases/%s", ApiURL, ApiVersion, id.String())
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

type DBQueryResponse struct {
	Object     ObjectType   `json:"object"`
	Results    []PageObject `json:"results"`
	HasMore    bool         `json:"has_more"`
	NextCursor Cursor       `json:"next_cursor"`
}
