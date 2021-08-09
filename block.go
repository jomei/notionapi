package notionapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BlockID string

func (bID BlockID) String() string {
	return string(bID)
}

type BlockService interface {
	GetChildren(context.Context, BlockID, *Pagination) (*GetChildrenResponse, error)
	AppendChildren(context.Context, BlockID, *AppendBlockChildrenRequest) (Block, error)
}

type BlockClient struct {
	apiClient *Client
}

// GetChildren https://developers.notion.com/reference/get-block-children
func (bc *BlockClient) GetChildren(ctx context.Context, id BlockID, pagination *Pagination) (*GetChildrenResponse, error) {
	res, err := bc.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("blocks/%s/children", id.String()), pagination.ToQuery(), nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Object  ObjectType `json:"object"`
		Results []map[string]interface{}
	}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	results := make([]Block, len(response.Results))
	for i, r := range response.Results {
		b, err := decodeBlock(r)
		if err != nil {
			return nil, err
		}
		results[i] = b
	}

	return &GetChildrenResponse{
		Object:  response.Object,
		Results: results,
	}, nil
}

type GetChildrenResponse struct {
	Object  ObjectType `json:"object"`
	Results []Block    `json:"results"`
}

// AppendChildren https://developers.notion.com/reference/patch-block-children
func (bc *BlockClient) AppendChildren(ctx context.Context, id BlockID, requestBody *AppendBlockChildrenRequest) (Block, error) {
	res, err := bc.apiClient.request(ctx, http.MethodPatch, fmt.Sprintf("blocks/%s/children", id.String()), nil, requestBody)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return decodeBlock(response)
}

type BlockType string

func (bt BlockType) String() string {
	return string(bt)
}

type AppendBlockChildrenRequest struct {
	Children []Block `json:"children"`
}

type Block interface {
	GetType() BlockType
}

type ParagraphBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Paragraph      struct {
		Text     Paragraph `json:"text"`
		Children []Block   `json:"children,omitempty"`
	} `json:"paragraph"`
}

func (b *ParagraphBlock) GetType() BlockType {
	return b.Type
}

type Heading1Block struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Heading1       struct {
		Text Paragraph `json:"text"`
	} `json:"heading_1"`
}

func (b *Heading1Block) GetType() BlockType {
	return b.Type
}

type Heading2Block struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Heading2       struct {
		Text Paragraph `json:"text"`
	} `json:"heading_2"`
}

func (b *Heading2Block) GetType() BlockType {
	return b.Type
}

type Heading3Block struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Heading3       struct {
		Text Paragraph `json:"text"`
	} `json:"heading_3"`
}

func (b *Heading3Block) GetType() BlockType {
	return b.Type
}

type BulletedListItemBlock struct {
	Object           ObjectType `json:"object"`
	ID               BlockID    `json:"id,omitempty"`
	Type             BlockType  `json:"type"`
	CreatedTime      *time.Time `json:"created_time,omitempty"`
	LastEditedTime   *time.Time `json:"last_edited_time,omitempty"`
	HasChildren      bool       `json:"has_children,omitempty"`
	BulletedListItem struct {
		Text     Paragraph `json:"text"`
		Children []Block   `json:"children,omitempty"`
	} `json:"bulleted_list_item"`
}

func (b *BulletedListItemBlock) GetType() BlockType {
	return b.Type
}

type NumberedListItemBlock struct {
	Object           ObjectType `json:"object"`
	ID               BlockID    `json:"id,omitempty"`
	Type             BlockType  `json:"type"`
	CreatedTime      *time.Time `json:"created_time,omitempty"`
	LastEditedTime   *time.Time `json:"last_edited_time,omitempty"`
	HasChildren      bool       `json:"has_children,omitempty"`
	NumberedListItem struct {
		Text     Paragraph `json:"text"`
		Children []Block   `json:"children,omitempty"`
	} `json:"numbered_list_item"`
}

func (b *NumberedListItemBlock) GetType() BlockType {
	return b.Type
}

type ToDoBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children"`
	ToDo           struct {
		Text     Paragraph `json:"text"`
		Children []Block   `json:"children,omitempty"`
		Checked  bool      `json:"checked"`
	} `json:"to_do"`
}

func (b *ToDoBlock) GetType() BlockType {
	return b.Type
}

type ToggleBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Text           Paragraph  `json:"text"`
	Children       []Block    `json:"children,omitempty"`
	Toggle         struct {
		Text     Paragraph `json:"text"`
		Children []Block   `json:"children,omitempty"`
	} `json:"toggle"`
}

func (b *ToggleBlock) GetType() BlockType {
	return b.Type
}

type ChildPageBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	ChildPage      struct {
		Title string `json:"title"`
	} `json:"child_page"`
}

func (b *ChildPageBlock) GetType() BlockType {
	return b.Type
}

func decodeBlock(raw map[string]interface{}) (Block, error) {
	var b Block
	switch BlockType(raw["type"].(string)) {
	case BlockTypeParagraph:
		b = &ParagraphBlock{}
	case BlockTypeHeading1:
		b = &Heading1Block{}
	case BlockTypeHeading2:
		b = &Heading2Block{}
	case BlockTypeHeading3:
		b = &Heading2Block{}
	case BlockTypeBulletedListItem:
		b = &BulletedListItemBlock{}
	case BlockTypeNumberedListItem:
		b = &NumberedListItemBlock{}
	case BlockTypeToDo:
		b = &ToDoBlock{}
	case BlockTypeToggle:
		b = &ToggleBlock{}
	case BlockTypeChildPage:
		b = &ChildPageBlock{}
	default:
		return nil, fmt.Errorf("unsupported block type: %s", raw["type"].(string))
	}
	j, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(j, b)
	return b, err
}
