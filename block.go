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
	AppendChildren(context.Context, BlockID, *AppendBlockChildrenRequest) (*AppendBlockChildrenResponse, error)
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
		Object     ObjectType `json:"object"`
		NextCursor string     `json:"next_cursor"`
		HasMore    bool       `json:"has_more"`
		Results    []map[string]interface{}
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
		Object:     response.Object,
		Results:    results,
		NextCursor: response.NextCursor,
		HasMore:    response.HasMore,
	}, nil
}

type GetChildrenResponse struct {
	Object     ObjectType `json:"object"`
	Results    []Block    `json:"results"`
	NextCursor string     `json:"next_cursor"`
	HasMore    bool       `json:"has_more"`
}

// AppendChildren https://developers.notion.com/reference/patch-block-children
func (bc *BlockClient) AppendChildren(ctx context.Context, id BlockID, requestBody *AppendBlockChildrenRequest) (*AppendBlockChildrenResponse, error) {
	res, err := bc.apiClient.request(ctx, http.MethodPatch, fmt.Sprintf("blocks/%s/children", id.String()), nil, requestBody)
	if err != nil {
		return nil, err
	}

	var response AppendBlockChildrenResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
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

type CalloutBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Callout        Callout    `json:"callout"`
}

func (b CalloutBlock) GetType() BlockType {
	return b.Type
}

type Callout struct {
	Text     []RichText `json:"text"`
	Icon     *Icon      `json:"icon,omitempty"`
	Children []Block    `json:"children,omitempty"`
}

type QuoteBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Quote          Quote      `json:"quote"`
}

func (b QuoteBlock) GetType() BlockType {
	return b.Type
}

type Quote struct {
	Text     []RichText `json:"text"`
	Children []Block    `json:"children,omitempty"`
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
	Caption  []RichText  `json:"caption,omitempty"`
	Type     FileType    `json:"type"`
	File     *FileObject `json:"file,omitempty"`
	External *FileObject `json:"external,omitempty"`
}

// GetURL returns the external or internal URL depending on the image type.
func (i Image) GetURL() string {
	if i.File != nil {
		return i.File.URL
	}
	if i.External != nil {
		return i.External.URL
	}
	return ""
}

type CodeBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	Code           Code       `json:"code"`
}

func (b CodeBlock) GetType() BlockType {
	return b.Type
}

type Code struct {
	Text     []RichText `json:"text"`
	Language string     `json:"language"`
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
	Caption  []RichText  `json:"caption,omitempty"`
	Type     FileType    `json:"type"`
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
	Caption  []RichText  `json:"caption,omitempty"`
	Type     FileType    `json:"type"`
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
	Caption  []RichText  `json:"caption,omitempty"`
	Type     FileType    `json:"type"`
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

type ChildDatabaseBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	ChildDatabase  struct {
		Title string `json:"title"`
	} `json:"child_database"`
}

func (b ChildDatabaseBlock) GetType() BlockType {
	return b.Type
}

type TableOfContentsBlock struct {
	Object         ObjectType     `json:"object"`
	ID             BlockID        `json:"id,omitempty"`
	Type           BlockType      `json:"type"`
	CreatedTime    *time.Time     `json:"created_time,omitempty"`
	LastEditedTime *time.Time     `json:"last_edited_time,omitempty"`
	HasChildren    bool           `json:"has_children,omitempty"`
	TableOfContent TableOfContent `json:"table_of_contents"`
}

func (b TableOfContentsBlock) GetType() BlockType {
	return b.Type
}

type TableOfContent struct {
	// empty
}

type DividerBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`
	Divider        Divider    `json:"divider"`
}

func (b DividerBlock) GetType() BlockType {
	return b.Type
}

type Divider struct {
	// empty
}

type UnsupportedBlock struct {
	Object         ObjectType `json:"object"`
	ID             BlockID    `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
}

func (b UnsupportedBlock) GetType() BlockType {
	return b.Type
}

type AppendBlockChildrenResponse struct {
	Object  ObjectType `json:"object"`
	Results []Block    `json:"results"`
}

type appendBlockResponse struct {
	Object  ObjectType               `json:"object"`
	Results []map[string]interface{} `json:"results"`
}

func (r *AppendBlockChildrenResponse) UnmarshalJSON(data []byte) error {
	var raw appendBlockResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	blocks := make([]Block, 0)
	for _, b := range raw.Results {
		block, err := decodeBlock(b)
		if err != nil {
			return err
		}
		blocks = append(blocks, block)
	}

	*r = AppendBlockChildrenResponse{
		Object:  raw.Object,
		Results: blocks,
	}
	return nil
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
		b = &Heading3Block{}
	case BlockCallout:
		b = &CalloutBlock{}
	case BlockQuote:
		b = &QuoteBlock{}
	case BlockTypeBulletedListItem:
		b = &BulletedListItemBlock{}
	case BlockTypeNumberedListItem:
		b = &NumberedListItemBlock{}
	case BlockTypeToDo:
		b = &ToDoBlock{}
	case BlockTypeCode:
		b = &CodeBlock{}
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
	case BlockTypeChildDatabase:
		b = &ChildDatabaseBlock{}
	case BlockTypeTableOfContents:
		b = &TableOfContentsBlock{}
	case BlockTypeDivider:
		b = &DividerBlock{}
	case BlockTypeUnsupported:
		b = &UnsupportedBlock{}
	default:
		return &UnsupportedBlock{}, nil
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
	Toggle           *Toggle    `json:"toggle,omitempty"`
	Embed            *Embed     `json:"embed,omitempty"`
	Image            *Image     `json:"image,omitempty"`
	Video            *Video     `json:"video,omitempty"`
	File             *BlockFile `json:"file,omitempty"`
	Pdf              *Pdf       `json:"pdf,omitempty"`
	Bookmark         *Bookmark  `json:"bookmark,omitempty"`
}
