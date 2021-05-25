package notionapi

import "time"

type PageID string

func (pID PageID) String() string {
	return string(pID)
}

type PageObject struct {
	Object         ObjectType              `json:"object"`
	ID             ObjectID                `json:"id"`
	CreatedTime    time.Time               `json:"created_time"` // TODO: format
	LastEditedTime time.Time               `json:"last_edited_time"`
	Archived       bool                    `json:"archived"`
	Properties     map[PropertyName]Object `json:"properties"`
	Parent         Parent                  `json:"parent"`
}

type Parent struct {
	Type       ObjectType `json:"type"`
	PageID     PageID     `json:"page_id,omitempty"`
	DatabaseID DatabaseID `json:"database_id,omitempty"`
}
