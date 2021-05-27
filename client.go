package notionapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
)

const (
	apiURL        = "https://api.notion.com"
	apiVersion    = "v1"
	notionVersion = "2021-05-13"
)

type Token string

func (it Token) String() string {
	return string(it)
}

// ClientOption to configure API client
type ClientOption func(*Client)

type Client struct {
	httpClient    *http.Client
	baseUrl       *url.URL
	apiVersion    string
	notionVersion string

	Token Token

	Database DatabaseService
	Block    BlockService
	Page     PageService
	User     UserService
	Search   SearchService
}

func NewClient(token Token, opts ...ClientOption) *Client {
	u, err := url.Parse(apiURL)
	if err != nil {
		panic(err)
	}
	c := &Client{
		httpClient:    http.DefaultClient,
		Token:         token,
		baseUrl:       u,
		apiVersion:    apiVersion,
		notionVersion: notionVersion,
	}

	c.Database = &DatabaseClient{apiClient: c}
	c.Block = &BlockClient{apiClient: c}
	c.Page = &PageClient{apiClient: c}
	c.User = &UserClient{apiClient: c}
	c.Search = &SearchClient{apiClient: c}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithHTTPClient overrides the default http.Client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithVersion overrides the Notion API version
func WithVersion(version string) ClientOption {
	return func(c *Client) {
		c.notionVersion = version
	}
}

func (c *Client) request(ctx context.Context, method string, urlStr string, queryParams map[string]string, requestBody interface{}) (*http.Response, error) {
	u, err := c.baseUrl.Parse(fmt.Sprintf("%s/%s", c.apiVersion, urlStr))
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if requestBody != nil {
		body, err := json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(body)
	}

	if len(queryParams) > 0 {
		q := u.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token.String()))
	req.Header.Add("Notion-Version", c.notionVersion)

	if requestBody != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("http status: %d", res.StatusCode)
	}

	return res, nil
}
