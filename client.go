package notionapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const (
	apiURL        = "https://api.notion.com"
	apiVersion    = "v1"
	notionVersion = "2022-06-28"
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

	requestMutex sync.Mutex

	Token Token

	Database DatabaseService
	Block    BlockService
	Page     PageService
	User     UserService
	Search   SearchService
	Comment  CommentService
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
	c.Comment = &CommentClient{apiClient: c}

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
	c.requestMutex.Lock()
	defer c.requestMutex.Unlock()

	u, err := c.baseUrl.Parse(fmt.Sprintf("%s/%s", c.apiVersion, urlStr))
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if requestBody != nil && !reflect.ValueOf(requestBody).IsNil() {
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
	req.Header.Add("Content-Type", "application/json")

	failedAttempts := 0
	var res *http.Response
	for {
		var err error
		res, err = c.httpClient.Do(req.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		if res.StatusCode == http.StatusTooManyRequests {
			failedAttempts++
			if failedAttempts == 3 {
				return nil, errors.New("Retry request with 429 response failed")
			}
			// https://developers.notion.com/reference/request-limits#rate-limits
			retryAfter := res.Header["Retry-After"][0]
			waitSeconds, err := strconv.Atoi(retryAfter)
			if err != nil {
				break // should not happen
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(waitSeconds) * time.Second):
			}
		} else {
			break
		}
	}

	if res.StatusCode != http.StatusOK {
		var apiErr Error
		err = json.NewDecoder(res.Body).Decode(&apiErr)
		if err != nil {
			return nil, err
		}

		return nil, &apiErr
	}

	return res, nil
}

type Pagination struct {
	StartCursor Cursor
	PageSize    int
}

func (p *Pagination) ToQuery() map[string]string {
	if p == nil {
		return nil
	}
	r := map[string]string{}
	if p.StartCursor != "" {
		r["start_cursor"] = p.StartCursor.String()
	}

	if p.PageSize != 0 {
		r["page_size"] = strconv.Itoa(p.PageSize)
	}

	return r
}
