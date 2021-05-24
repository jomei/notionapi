package notionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type DatabaseID string

func (dID DatabaseID) String() string {
	return string(dID)
}

func (c *Client) DBRetrieve(ctx context.Context, id DatabaseID) (*DBRetrieveResponse, error) {
	req, err := makeDBRetrieveRequest(id, c.Token)
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

func makeDBRetrieveRequest(id DatabaseID, token IntegrationToken) (*http.Request, error) {
	reqURL := fmt.Sprintf("%s/%s/databases/%s", ApiURL, ApiVersion, id.String())
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("application/json", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.String()))
	req.Header.Add("Notion-Version", NotionVersion)

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
