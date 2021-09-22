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
	Get(context.Context, BlockID) (Block, error)
	Update(ctx context.Context, id BlockID, request *BlockUpdateRequest) (Block, error)
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

// Get https://developers.notion.com/reference/retrieve-a-block
func (bc *BlockClient) Get(ctx context.Context, id BlockID) (Block, error) {
	res, err := bc.apiClient.request(ctx, http.MethodGet, fmt.Sprintf("blocks/%s", id.String()), nil, nil)
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

func (bc *BlockClient) Update(ctx context.Context, id BlockID, requestBody *BlockUpdateRequest) (Block, error) {
	res, err := bc.apiClient.request(ctx, http.MethodPatch, fmt.Sprintf("blocks/%s", id.String()), nil, requestBody)
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
	Paragraph      Paragraph  `json:"paragraph"`
}

type Paragraph struct {
	Text     []RichText `json:"text"`
	Children []Block    `json:"children,omitempty"`
}

func (b ParagraphBlock) GetType() BlockType {
	return b.Type
}

type Heading1Block struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Heading1       Heading    `json:"heading_1"`
}

type Heading struct {
	Text []RichText `json:"text"`
}

func (b Heading1Block) GetType() BlockType {
	return b.Type
}

type Heading2Block struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Heading2       Heading    `json:"heading_2"`
}

func (b Heading2Block) GetType() BlockType {
	return b.Type
}

type Heading3Block struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Heading3       Heading    `json:"heading_3"`
}

func (b Heading3Block) GetType() BlockType {
	return b.Type
}

type BulletedListItemBlock struct {
	Object           ObjectType `json:"object"`
	ID               BlockID    `json:"id,omitempty"`
	Type             BlockType  `json:"type"`
	CreatedTime      *time.Time `json:"created_time,omitempty"`
	LastEditedTime   *time.Time `json:"last_edited_time,omitempty"`
	HasChildren      bool       `json:"has_children,omitempty"`
	BulletedListItem ListItem   `json:"bulleted_list_item"`
}

type ListItem struct {
	Text     []RichText `json:"text"`
	Children []Block    `json:"children,omitempty"`
}

func (b BulletedListItemBlock) GetType() BlockType {
	return b.Type
}

type NumberedListItemBlock struct {
	Object           ObjectType `json:"object"`
	ID               BlockID    `json:"id,omitempty"`
	Type             BlockType  `json:"type"`
	CreatedTime      *time.Time `json:"created_time,omitempty"`
	LastEditedTime   *time.Time `json:"last_edited_time,omitempty"`
	HasChildren      bool       `json:"has_children,omitempty"`
	NumberedListItem ListItem   `json:"numbered_list_item"`
}

func (b NumberedListItemBlock) GetType() BlockType {
	return b.Type
}

type ToDoBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children"`
	ToDo           ToDo       `json:"to_do"`
}

type ToDo struct {
	Text     []RichText `json:"text"`
	Children []Block    `json:"children,omitempty"`
	Checked  bool       `json:"checked,omitempty"`
}

func (b ToDoBlock) GetType() BlockType {
	return b.Type
}

type ToggleBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Text           []RichText `json:"text"`
	Children       []Block    `json:"children,omitempty"`
	Toggle         Toggle     `json:"toggle"`
}

type Toggle struct {
	Text     []RichText `json:"text"`
	Children []Block    `json:"children,omitempty"`
}

func (b ToggleBlock) GetType() BlockType {
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

func (b ChildPageBlock) GetType() BlockType {
	return b.Type
}

type EmbedBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Embed          Embed      `json:"embed"`
}

func (b EmbedBlock) GetType() BlockType {
	return b.Type
}

type Embed struct {
	Caption []RichText `json:"caption,omitempty"`
	URL     string     `json:"url"`
}

type ImageBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Image          Image      `json:"image"`
}

func (b ImageBlock) GetType() BlockType {
	return b.Type
}

type Image struct {
	Caption  []RichText `json:"caption,omitempty"`
	Type     FileType   `json:"type"`
	File     *FileObject `json:"file,omitempty"`
	External *FileObject `json:"external,omitempty"`
}

type VideoBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Video          Video      `json:"video"`
}

func (b VideoBlock) GetType() BlockType {
	return b.Type
}

type Video struct {
	Caption  []RichText `json:"caption,omitempty"`
	Type     FileType   `json:"type"`
	File     *FileObject `json:"file,omitempty"`
	External *FileObject `json:"external,omitempty"`
}

type FileBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	File           BlockFile  `json:"file"`
}

func (b FileBlock) GetType() BlockType {
	return b.Type
}

type BlockFile struct {
	Caption  []RichText `json:"caption,omitempty"`
	Type     FileType   `json:"type"`
	File     *FileObject `json:"file,omitempty"`
	External *FileObject `json:"external,omitempty"`
}

type PdfBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Pdf            Pdf        `json:"pdf"`
}

func (b PdfBlock) GetType() BlockType {
	return b.Type
}

type Pdf struct {
	Caption  []RichText `json:"caption,omitempty"`
	Type     FileType   `json:"type"`
	File     *FileObject `json:"file,omitempty"`
	External *FileObject `json:"external,omitempty"`
}

type BookmarkBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Bookmark       Bookmark   `json:"bookmark"`
}

func (b BookmarkBlock) GetType() BlockType {
	return b.Type
}

type Bookmark struct {
	Caption []RichText `json:"caption,omitempty"`
	URL     string     `json:"url"`
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
	case BlockTypeEmbed:
		b = &EmbedBlock{}
	case BlockTypeImage:
		b = &ImageBlock{}
	case BlockTypeVideo:
		b = &VideoBlock{}
	case BlockTypeFile:
		b = &FileBlock{}
	case BlockTypePdf:
		b = &PdfBlock{}
	case BlockTypeBookmark:
		b = &BookmarkBlock{}
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

type BlockUpdateRequest struct {
	Paragraph        *Paragraph `json:"paragraph,omitempty"`
	Heading1         *Heading   `json:"heading_1,omitempty"`
	Heading2         *Heading   `json:"heading_2,omitempty"`
	Heading3         *Heading   `json:"heading_3,omitempty"`
	BulletedListItem *ListItem  `json:"bulleted_list_item,omitempty"`
	NumberedListItem *ListItem  `json:"numbered_list_item,omitempty"`
	ToDo             *ToDo      `json:"to_do,omitempty"`
	Toggle           *Toggle    `json:"toggle,omtiempty"`
	Embed            *Embed     `json:"embed,omitempty"`
	Image            *Image     `json:"image,omitempty"`
	Video            *Video     `json:"video,omitempty"`
	File             *BlockFile `json:"file,omitempty"`
	Pdf              *Pdf       `json:"pdf,omitempty"`
	Bookmark         *Bookmark  `json:"bookmark,omitempty"`
}
