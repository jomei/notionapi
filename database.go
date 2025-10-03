package notionapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type DatabaseID string

func (dID DatabaseID) String() string {
	return string(dID)
}

type DatabaseService interface {
	Create(ctx context.Context, request *DatabaseCreateRequest) (*Database, error)
	Get(context.Context, DatabaseID) (*Database, error)
	Update(context.Context, DatabaseID, *DatabaseUpdateRequest) (*Database, error)
}

type DatabaseClient struct {
	apiClient *Client
}

// Creates a database as a subpage in the specified parent page, with the
// specified properties schema. Currently, the parent of a new database must be
// a Notion page.
//
// See https://developers.notion.com/reference/create-a-database
func (dc *DatabaseClient) Create(ctx context.Context, requestBody *DatabaseCreateRequest) (*Database, error) {
	res, err := dc.apiClient.request(ctx, http.MethodPost, "databases", nil, requestBody)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	var response Database
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DatabaseCreateRequest represents the request body for DatabaseClient.Create.
type DatabaseCreateRequest struct {
	// A page parent.
	Parent Parent `json:"parent"`
	// Title of database as it appears in Notion. An array of rich text objects.
	Title []RichText `json:"title"`
	// Property schema of database. The keys are the names of properties as they
	// appear in Notion and the values are property schema objects.
	Properties PropertyConfigs `json:"properties"`
	IsInline   bool            `json:"is_inline"`
}

// See https://developers.notion.com/reference/get-database
func (dc *DatabaseClient) Get(ctx context.Context, id DatabaseID) (*Database, error) {
	if id == "" {
		return nil, errors.New("empty database id")
	}

	res, err := dc.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("databases/%s", id.String()), nil, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	var response Database

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Update https://developers.notion.com/reference/update-a-database
func (dc *DatabaseClient) Update(ctx context.Context, id DatabaseID, requestBody *DatabaseUpdateRequest) (*Database, error) {
	res, err := dc.apiClient.request(ctx, http.MethodPatch, fmt.Sprintf("databases/%s", id.String()), nil, requestBody)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := res.Body.Close(); errClose != nil {
			log.Println("failed to close body, should never happen")
		}
	}()

	var response Database
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

// DatabaseUpdateRequest represents the request body for DatabaseClient.Update.
type DatabaseUpdateRequest struct {
	// An array of rich text objects that represents the title of the database
	// that is displayed in the Notion UI. If omitted, then the database title
	// remains unchanged.
	Title []RichText `json:"title,omitempty"`
	// The properties of a database to be changed in the request, in the form of
	// a JSON object. If updating an existing property, then the keys are the
	// names or IDs of the properties as they appear in Notion, and the values are
	// property schema objects. If adding a new property, then the key is the name
	// of the new database property and the value is a property schema object.
	Properties PropertyConfigs `json:"properties,omitempty"`
}

type Database struct {
	Object         ObjectType `json:"object"`
	ID             ObjectID   `json:"id"`
	CreatedTime    time.Time  `json:"created_time"`
	LastEditedTime time.Time  `json:"last_edited_time"`
	CreatedBy      User       `json:"created_by,omitempty"`
	LastEditedBy   User       `json:"last_edited_by,omitempty"`
	Title          []RichText `json:"title"`
	Parent         Parent     `json:"parent"`
	URL            string     `json:"url"`
	PublicURL      string     `json:"public_url"`
	// Properties is a map of property configurations that defines what Page.Properties each page of the database can use
	Properties  PropertyConfigs      `json:"properties"`
	Description []RichText           `json:"description"`
	IsInline    bool                 `json:"is_inline"`
	Archived    bool                 `json:"archived"`
	Icon        *Icon                `json:"icon,omitempty"`
	Cover       *Image               `json:"cover,omitempty"`
	DataSources []DatabaseDataSource `json:"data_sources"`
}

type DatabaseDataSource struct {
	ID   DataSourceID `json:"id"`
	Name string       `json:"name"`
}

func (db *Database) GetObject() ObjectType {
	return db.Object
}
