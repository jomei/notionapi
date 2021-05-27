package notionapi

type PropertyType string

type Property interface {
	GetType() PropertyType
}

type PropertyID string

func (pID PropertyID) String() string {
	return string(pID)
}

type TextProperty struct {
	ID    PropertyID   `json:"id"`
	Type  PropertyType `json:"type"`
	Title []TextObject `json:"title"`
}

func (p TextProperty) GetType() PropertyType {
	return p.Type
}

type RichTextProperty struct {
	ID       PropertyID   `json:"id"`
	Type     PropertyType `json:"type"`
	RichText TextObject   `json:"rich_text"`
}

func (p RichTextProperty) GetType() PropertyType {
	return p.Type
}

type TitleProperty struct {
	ID    PropertyID   `json:"id"`
	Type  PropertyType `json:"type"`
	Title TextObject   `json:"rich_text"`
}

func (p TitleProperty) GetType() PropertyType {
	return p.Type
}

type FormatType string

func (ft FormatType) String() string {
	return string(ft)
}

type NumberProperty struct {
	ID     ObjectID     `json:"id"`
	Type   PropertyType `json:"type"`
	Format FormatType   `json:"format"`
}

func (p NumberProperty) GetType() PropertyType {
	return p.Type
}

type SelectProperty struct {
	ID     ObjectID     `json:"id"`
	Type   PropertyType `json:"type"`
	Select Select       `json:"select"`
}

type Select struct {
	Options []Option `json:"options"`
}

type MultiSelectProperty struct {
	ID          ObjectID     `json:"id"`
	Type        PropertyType `json:"type"`
	MultiSelect Select       `json:"multi_select"`
}

func (p MultiSelectProperty) GetType() PropertyType {
	return p.Type
}

type Option struct {
	ID    PropertyID
	Name  string `json:"name"`
	Color Color  `json:"color"`
}

func (p SelectProperty) GetType() PropertyType {
	return p.Type
}

type DateProperty struct {
	ID   ObjectID     `json:"id"`
	Type PropertyType `json:"type"`
	Date interface{}  `json:"date"`
}

func (p DateProperty) GetType() PropertyType {
	return p.Type
}

type PeopleProperty struct {
	ID     ObjectID     `json:"id"`
	Type   PropertyType `json:"type"`
	People interface{}  `json:"people"`
}

func (p PeopleProperty) GetType() PropertyType {
	return p.Type
}

type FileProperty struct {
	ID   ObjectID     `json:"id"`
	Type PropertyType `json:"type"`
	File interface{}  `json:"file"`
}

func (p FileProperty) GetType() PropertyType {
	return p.Type
}

type CheckboxProperty struct {
	ID       ObjectID     `json:"id"`
	Type     PropertyType `json:"type"`
	Checkbox interface{}  `json:"checkbox"`
}

func (p CheckboxProperty) GetType() PropertyType {
	return p.Type
}

type URLProperty struct {
	ID   ObjectID     `json:"id"`
	Type PropertyType `json:"type"`
	URL  interface{}  `json:"url"`
}

func (p URLProperty) GetType() PropertyType {
	return p.Type
}

type EmailProperty struct {
	ID    PropertyID   `json:"id"`
	Type  PropertyType `json:"type"`
	Email interface{}  `json:"email"`
}

func (p EmailProperty) GetType() PropertyType {
	return p.Type
}

type PhoneNumberProperty struct {
	ID          ObjectID     `json:"id"`
	Type        PropertyType `json:"type"`
	PhoneNumber interface{}  `json:"phone_number"`
}

func (p PhoneNumberProperty) GetType() PropertyType {
	return p.Type
}

type FormulaProperty struct {
	ID         ObjectID     `json:"id"`
	Type       PropertyType `json:"type"`
	Expression string       `json:"expression"`
}

func (p FormulaProperty) GetType() PropertyType {
	return p.Type
}

type RelationProperty struct {
	Type     PropertyType `json:"type"`
	Relation Relation     `json:"relation"`
}

type Relation struct {
	DatabaseID         DatabaseID `json:"database_id"`
	SyncedPropertyID   PropertyID `json:"synced_property_id"`
	SyncedPropertyName string     `json:"synced_property_name"`
}

func (p RelationProperty) GetType() PropertyType {
	return p.Type
}

type RollupProperty struct {
	ID     ObjectID     `json:"id"`
	Type   PropertyType `json:"type"`
	Rollup Rollup       `json:"rollup"`
}

type Rollup struct {
	RelationPropertyName string       `json:"relation_property_name"`
	RelationPropertyID   PropertyID   `json:"relation_property_id"`
	RollupPropertyName   string       `json:"rollup_property_name"`
	RollupPropertyID     PropertyID   `json:"rollup_property_id"`
	Function             FunctionType `json:"function"`
}

func (p RollupProperty) GetType() PropertyType {
	return p.Type
}

type CreatedTimeProperty struct {
	ID          ObjectID     `json:"id"`
	Type        PropertyType `json:"type"`
	CreatedTime interface{}  `json:"created_time"`
}

func (p CreatedTimeProperty) GetType() PropertyType {
	return p.Type
}

type CreatedByProperty struct {
	ID        ObjectID     `json:"id"`
	Type      PropertyType `json:"type"`
	CreatedBy interface{}  `json:"created_by"`
}

func (p CreatedByProperty) GetType() PropertyType {
	return p.Type
}

type LastEditedTimeProperty struct {
	ID             ObjectID     `json:"id"`
	Type           PropertyType `json:"type"`
	LastEditedTime interface{}  `json:"last_edited_time"`
}

func (p LastEditedTimeProperty) GetType() PropertyType {
	return p.Type
}

type LastEditedByProperty struct {
	ID           ObjectID     `json:"id"`
	Type         PropertyType `json:"type"`
	LastEditedBy interface{}  `json:"last_edited_by"`
}

func (p LastEditedByProperty) GetType() PropertyType {
	return p.Type
}
