package notionapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

const (
	apiURL        = "https://api.notion.com"
	apiVersion    = "v1"
	notionVersion = "2021-05-13"
)

type ApiClient interface {
	DBRetrieve(context.Context, DatabaseID) (*DatabaseObject, error)
	DBList(context.Context, Cursor, int) (*DBListResponse, error)
	DBQuery(context.Context, DatabaseID, DBQueryRequest) (*DBQueryResponse, error)
	BlockChildrenRetrieve(context.Context, BlockID, Cursor, int) ([]BasicObject, error)
	BlockChildrenAppend(context.Context, BlockID, AppendBlockChildrenRequest) (*BlockObject, error)
	PageRetrieve(context.Context, PageID) (*PageObject, error)
	PageCreate(context.Context, PageCreateRequest) (*PageObject, error)
	PageUpdate(context.Context, PageID, map[PropertyName]BasicObject) (*PageObject, error)
	UserRetrieve(context.Context, UserID) (*UserObject, error)
	UsersList(context.Context, Cursor, int) (*UsersListResponse, error)
	Search(context.Context, SearchRequest) (*SearchResponse, error)
}

// ClientOption to configure API client
type ClientOption func(*Client)

type Client struct {
	httpClient *http.Client

	Token         Token
	baseUrl       *url.URL
	apiVersion    string
	notionVersion string
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

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type Token string

func (it Token) String() string {
	return string(it)
}

func (c *Client) request(ctx context.Context, method string, urlStr string, queryParams map[string]string, requestBody interface{}) (*http.Response, error) {
	u, err := c.baseUrl.Parse(fmt.Sprintf("%s/%s", c.apiVersion, urlStr))
	if err != nil {
		return nil, err
	}

	var buf *bytes.Buffer
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

type Cursor string

func (c Cursor) String() string {
	return string(c)
}
