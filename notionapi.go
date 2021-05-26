package notionapi

import (
	"context"
	"fmt"
	"net/http"
)

const (
	ApiURL        = "https://api.notion.com"
	ApiVersion    = "v1"
	NotionVersion = "2021-05-13"
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

type ClientOption func(*Client)

type Client struct {
	httpClient *http.Client

	Token IntegrationToken
}

func NewClient(httpClient *http.Client, token IntegrationToken, opts ...ClientOption) *Client {
	c := &Client{
		httpClient: httpClient,
		Token:      token,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type IntegrationToken string

func (it IntegrationToken) String() string {
	return string(it)
}

type Color string

func (c Color) String() string {
	return string(c)
}

func (c *Client) addRequestHeaders(req *http.Request) *http.Request {
	req.Header.Add("application/json", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token.String()))
	req.Header.Add("Notion-Version", NotionVersion)

	return req
}

type Cursor string

func (c Cursor) String() string {
	return string(c)
}
