package notionapi

const (
	ObjectTypeTitle       ObjectType = "title"
	ObjectTypeRichText    ObjectType = "rich_text"
	ObjectTypeCheckbox    ObjectType = "checkbox"
	ObjectTypeSelect      ObjectType = "select"
	ObjectTypeNumber      ObjectType = "number"
	ObjectTypeFormula     ObjectType = "formula"
	ObjectTypeDate        ObjectType = "date"
	ObjectTypeRelation    ObjectType = "relation"
	ObjectTypeRollup      ObjectType = "rollup"
	ObjectTypeMultiSelect ObjectType = "multi_select"
	ObjectTypePeople      ObjectType = "people"
	ObjectTypeFiles       ObjectType = "files"
	ObjectTypeHeading1    ObjectType = "heading_1"
	ObjectTypeHeading2    ObjectType = "heading_2"
	ObjectTypeHeading3    ObjectType = "heading_3"
	ObjectTypeParagraph   ObjectType = "paragraph"
	ObjectTypeToggle      ObjectType = "toggle"
	ObjectTypeUser        ObjectType = "user"

	ObjectTypeBulletedListItem ObjectType = "bulleted_list_item"
	ObjectTypeNumberedListItem ObjectType = "numbered_list_item"

	ObjectTypeToDo        ObjectType = "to_do"
	ObjectTypeChildPage   ObjectType = "child_page"
	ObjectTypeUnsupported ObjectType = "unsupported"
)

type ObjectType string

func (ot ObjectType) String() string {
	return string(ot)
}

type ObjectID string

func (oID ObjectID) String() string {
	return string(oID)
}

type Object interface {
	GetObject() ObjectType
}

//TODO: dont need anymore
type BasicObject struct {
	ID          ObjectID           `json:"id"`
	Type        ObjectType         `json:"type"`
	Title       *Text              `json:"title,omitempty"`
	Text        *Text              `json:"text,omitempty"`
	RichText    *RichText          `json:"rich_text,omitempty"`
	Checkbox    *struct{}          `json:"checkbox,omitempty"`
	Formula     *FormulaObject     `json:"formula,omitempty"`
	Date        *struct{}          `json:"date,omitempty"`
	Relation    *RelationObject    `json:"relation,omitempty"`
	Rollup      *RollupObject      `json:"rollup,omitempty"`
	MultiSelect *MultiSelectObject `json:"multi_select,omitempty"`
	People      *struct{}          `json:"people,omitempty"`
	Files       *struct{}          `json:"files,omitempty"`
	Paragraph   *Paragraph         `json:"paragraph,omitempty"`
	Toggle      *Toggle            `json:"toggle,omitempty"`
}

type Color string

func (c Color) String() string {
	return string(c)
}

type RichText struct {
	Type        ObjectType   `json:"type,omitempty"`
	Text        Text         `json:"text"`
	Annotations *Annotations `json:"annotations,omitempty"`
	PlainText   string       `json:"plain_text,omitempty"`
	Href        string       `json:"href,omitempty"`
}

type Text struct {
	Content string `json:"content"`
	Link    string `json:"link,omitempty"`
}

type Annotations struct {
	Bold          bool  `json:"bold"`
	Italic        bool  `json:"italic"`
	Strikethrough bool  `json:"strikethrough"`
	Underline     bool  `json:"underline"`
	Code          bool  `json:"code"`
	Color         Color `json:"color"`
}

type Paragraph []RichText

type FormulaObject struct {
	Value string `json:"value"`
}

type RelationObject struct {
	Database           DatabaseID `json:"database"`
	SyncedPropertyName string     `json:"synced_property_name"`
}

type FunctionType string

func (ft FunctionType) String() string {
	return string(ft)
}

type RollupObject struct {
	RollupPropertyName   string       `json:"rollup_property_name"`
	RelationPropertyName string       `json:"relation_property_name"`
	RollupPropertyID     ObjectID     `json:"rollup_property_id"`
	RelationPropertyID   ObjectID     `json:"relation_property_id"`
	Function             FunctionType `json:"function"`
}

type MultiSelectObject struct {
	Options []Option `json:"options"`
}

type Toggle struct {
	Text RichText `json:"text"`
}

type Cursor string

func (c Cursor) String() string {
	return string(c)
}
