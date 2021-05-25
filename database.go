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

type DatabaseID string

func (dID DatabaseID) String() string {
	return string(dID)
}

func (c *Client) DBRetrieve(ctx context.Context, id DatabaseID) (*DBRetrieveResponse, error) {
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

	var response DBRetrieveResponse
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

type DBRetrieveResponse struct {
	Object         ObjectType        `json:"object"`
	ID             ObjectID          `json:"id"`
	CreatedTime    time.Time         `json:"created_time"` //TODO: format
	LastEditedTime time.Time         `json:"last_edited_time"`
	Title          TextObject        `json:"title"`
	Properties     map[string]Object `json:"properties"`
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

type DatabaseObject struct {
	Object     ObjectType        `json:"object"`
	ID         ObjectID          `json:"id"`
	Title      string            `json:"title"`
	Properties map[string]Object `json:"properties"`
}
