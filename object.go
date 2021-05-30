package notionapi

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
