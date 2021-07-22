package notionapi

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type PropertyType string

type Property interface {
	GetType() PropertyType
}

type PropertyID string

func (pID PropertyID) String() string {
	return string(pID)
}

type TextProperty struct {
	ID    PropertyID   `json:"id,omitempty"`
	Type  PropertyType `json:"type"`
	Title []RichText   `json:"title"`
}

func (p TextProperty) GetType() PropertyType {
	return p.Type
}

type EmptyRichTextProperty struct {
	ID       PropertyID   `json:"id,omitempty"`
	Type     PropertyType `json:"type"`
	RichText struct{}     `json:"rich_text"`
}

func (p EmptyRichTextProperty) GetType() PropertyType {
	return p.Type
}

type RichTextProperty struct {
	ID       PropertyID   `json:"id,omitempty"`
	Type     PropertyType `json:"type"`
	RichText []RichText   `json:"rich_text"`
}

func (p RichTextProperty) GetType() PropertyType {
	return p.Type
}

type DatabaseTitleProperty struct {
	ID    PropertyID   `json:"id,omitempty"`
	Type  PropertyType `json:"type"`
	Title RichText     `json:"title"`
}

func (p DatabaseTitleProperty) GetType() PropertyType {
	return p.Type
}

type PageTitleProperty struct {
	ID    PropertyID   `json:"id,omitempty"`
	Type  PropertyType `json:"type,omitempty"`
	Title Paragraph    `json:"title"`
}

func (p PageTitleProperty) GetType() PropertyType {
	return p.Type
}

type FormatType string

func (ft FormatType) String() string {
	return string(ft)
}

type NumberProperty struct {
	ID     ObjectID     `json:"id,omitempty"`
	Type   PropertyType `json:"type"`
	Format FormatType   `json:"format"`
}

func (p NumberProperty) GetType() PropertyType {
	return p.Type
}

type SelectProperty struct {
	ID     ObjectID     `json:"id,omitempty"`
	Type   PropertyType `json:"type"`
	Select Select       `json:"select"`
}

func (p SelectProperty) GetType() PropertyType {
	return p.Type
}

type SelectOptionProperty struct {
	ID     ObjectID     `json:"id,omitempty"`
	Type   PropertyType `json:"type"`
	Select Option       `json:"select"`
}

func (p SelectOptionProperty) GetType() PropertyType {
	return p.Type
}

type Select struct {
	Options []Option `json:"options"`
}

type MultiSelectProperty struct {
	ID          ObjectID     `json:"id,omitempty"`
	Type        PropertyType `json:"type"`
	MultiSelect Select       `json:"multi_select"`
}

func (p MultiSelectProperty) GetType() PropertyType {
	return p.Type
}

type MultiSelectOptionsProperty struct {
	ID          ObjectID     `json:"id,omitempty"`
	Type        PropertyType `json:"type"`
	MultiSelect []Option     `json:"multi_select"`
}

func (p MultiSelectOptionsProperty) GetType() PropertyType {
	return p.Type
}

type Option struct {
	ID    PropertyID `json:"id,omitempty"`
	Name  string     `json:"name"`
	Color Color      `json:"color,omitempty"`
}

type DateProperty struct {
	ID   ObjectID     `json:"id,omitempty"`
	Type PropertyType `json:"type"`
	Date interface{}  `json:"date"`
}

func (p DateProperty) GetType() PropertyType {
	return p.Type
}

type PeopleProperty struct {
	ID     ObjectID     `json:"id,omitempty"`
	Type   PropertyType `json:"type"`
	People interface{}  `json:"people"`
}

func (p PeopleProperty) GetType() PropertyType {
	return p.Type
}

type FileProperty struct {
	ID   ObjectID     `json:"id,omitempty"`
	Type PropertyType `json:"type"`
	File interface{}  `json:"file"`
}

func (p FileProperty) GetType() PropertyType {
	return p.Type
}

type CheckboxProperty struct {
	ID       ObjectID     `json:"id,omitempty"`
	Type     PropertyType `json:"type"`
	Checkbox interface{}  `json:"checkbox"`
}

func (p CheckboxProperty) GetType() PropertyType {
	return p.Type
}

type URLProperty struct {
	ID   ObjectID     `json:"id,omitempty"`
	Type PropertyType `json:"type"`
	URL  interface{}  `json:"url"`
}

func (p URLProperty) GetType() PropertyType {
	return p.Type
}

type EmailProperty struct {
	ID    PropertyID   `json:"id,omitempty"`
	Type  PropertyType `json:"type"`
	Email interface{}  `json:"email"`
}

func (p EmailProperty) GetType() PropertyType {
	return p.Type
}

type PhoneNumberProperty struct {
	ID          ObjectID     `json:"id,omitempty"`
	Type        PropertyType `json:"type"`
	PhoneNumber interface{}  `json:"phone_number"`
}

func (p PhoneNumberProperty) GetType() PropertyType {
	return p.Type
}

type FormulaProperty struct {
	ID         ObjectID     `json:"id,omitempty"`
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
	ID     ObjectID     `json:"id,omitempty"`
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
	ID          ObjectID     `json:"id,omitempty"`
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

type Properties map[string]Property

func (p *Properties) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	props, err := parseProperties(raw)
	if err != nil {
		return err
	}

	*p = props
	return nil
}

func parseProperties(raw map[string]interface{}) (map[string]Property, error) {
	result := make(map[string]Property)
	for k, v := range raw {
		var p Property
		switch rawProperty := v.(type) {
		case map[string]interface{}:
			switch PropertyType(rawProperty["type"].(string)) {
			case PropertyTypeTitle:
				switch v.(map[string]interface{})["title"].(type) {
				case map[string]interface{}:
					p = &DatabaseTitleProperty{}
				default:
					p = &PageTitleProperty{}
				}
			case PropertyTypeRichText:
				switch v.(map[string]interface{})["rich_text"].(type) {
				case map[string]interface{}:
					p = &EmptyRichTextProperty{}
				default:
					p = &RichTextProperty{}
				}
			case PropertyTypeSelect:
				selectMap, ok := v.(map[string]interface{})["select"].(map[string]interface{})
				if !ok {
					return nil, errors.Errorf("an error occured while parsing property type: %s", rawProperty)
				}
				_, found := selectMap["options"]
				if found {
					p = &SelectProperty{}
				} else {
					p = &SelectOptionProperty{}
				}
			case PropertyTypeMultiSelect:
				switch v.(map[string]interface{})["multi_select"].(type) {
				case map[string]interface{}:
					p = &MultiSelectProperty{}
				default:
					p = &MultiSelectOptionsProperty{}
				}
			case PropertyTypeNumber:
				p = &NumberProperty{}
			case PropertyTypeCheckbox:
				p = &CheckboxProperty{}
			case PropertyTypeEmail:
				p = &EmailProperty{}
			case PropertyTypeURL:
				p = &URLProperty{}
			case PropertyTypeFile:
				p = &FileProperty{}
			case PropertyTypePhoneNumber:
				p = PhoneNumberProperty{}
			case PropertyTypeFormula:
				p = &FormulaProperty{}
			case PropertyTypeDate:
				p = &DateProperty{}
			case PropertyTypeRelation:
				p = &RelationProperty{}
			case PropertyTypeRollup:
				p = &RollupProperty{}
			case PropertyTypePeople:
				p = &PeopleProperty{}
			case PropertyTypeCreatedTime:
				p = &CreatedTimeProperty{}
			case PropertyTypeCreatedBy:
				p = &CreatedByProperty{}
			case PropertyTypeLastEditedTime:
				p = &LastEditedTimeProperty{}
			case PropertyTypeLastEditedBy:
				p = &LastEditedByProperty{}
			default:
				return nil, errors.New(fmt.Sprintf("unsupported property type: %s", rawProperty["type"].(string)))
			}
			b, err := json.Marshal(rawProperty)
			if err != nil {
				return nil, err
			}

			if err = json.Unmarshal(b, &p); err != nil {
				return nil, err
			}

			result[k] = p
		default:
			return nil, errors.New(fmt.Sprintf("unsupported property format %T", v))
		}
	}

	return result, nil
}
