package notionapi

const (
	ObjectTypeDatabase    ObjectType = "database"
	ObjectTypeTitle       ObjectType = "title"
	ObjectTypeText        ObjectType = "text"
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
	ObjectTypeList        ObjectType = "list"
)

type ObjectType string

func (ot ObjectType) String() string {
	return string(ot)
}

type ObjectID string

func (oID ObjectID) String() string {
	return string(oID)
}

type Object struct {
	ID          ObjectID           `json:"id"`
	Type        ObjectType         `json:"type"`
	Title       *struct{}          `json:"title,omitempty"`
	Text        *struct{}          `json:"text,omitempty"`
	Checkbox    *struct{}          `json:"checkbox,omitempty"`
	Select      *SelectObject      `json:"select,omitempty"`
	Number      *NumberObject      `json:"number,omitempty"`
	Formula     *FormulaObject     `json:"formula,omitempty"`
	Date        *struct{}          `json:"date,omitempty"`
	Relation    *RelationObject    `json:"relation,omitempty"`
	Rollup      *RollupObject      `json:"rollup,omitempty"`
	MultiSelect *MultiSelectObject `json:"multi_select,omitempty"`
	People      *struct{}          `json:"people,omitempty"`
	Files       *struct{}          `json:"files,omitempty"`
}

type TextObject struct {
	Type ObjectType `json:"type"`
	Text struct {
		Content string `json:"content"`
		Link    string `json:"link"`
	} `json:"text"`
	Annotations struct {
		Bold          bool  `json:"bold"`
		Italic        bool  `json:"italic"`
		Strikethrough bool  `json:"strikethrough"`
		Underline     bool  `json:"underline"`
		Code          bool  `json:"code"`
		Color         Color `json:"color"`
	} `json:"annotations"`
	PlainText string `json:"plain_text"`
	Href      string `json:"href"`
}

type SelectObject struct {
	Options []SelectObject `json:"options"`
}

type SelectOption struct {
	ID    ObjectID
	Name  string `json:"name"`
	Color Color  `json:"color"`
}

type FormatType string

func (ft FormatType) String() string {
	return string(ft)
}

type NumberObject struct {
	Format FormatType `json:"format"`
}

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
	Options [][]SelectOption `json:"options"`
}
