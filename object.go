package notionapi

import (
	"time"
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

type RelationObject struct {
	Database           DatabaseID `json:"database"`
	SyncedPropertyName string     `json:"synced_property_name"`
}

type FunctionType string

func (ft FunctionType) String() string {
	return string(ft)
}

type Cursor string

func (c Cursor) String() string {
	return string(c)
}

type Date time.Time

func (d *Date) String() string {
	return time.Time(*d).Format(time.RFC3339)
}

func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Date) UnmarshalText(data []byte) error {
	t, err := time.Parse(time.RFC3339, string(data))

	// Because the API does not distinguish between datetime with a
	// timezone and dates, we eventually have to try both.
	if err != nil {
		if _, ok := err.(*time.ParseError); !ok {
			return err
		} else {
			t, err = time.Parse("2006-01-02", string(data)) // Date
			if err != nil {
				// Still cannot parse it, nothing else to try.
				return err
			}
		}
	}

	*d = Date(t)
	return nil
}

type File struct {
	Name string `json:"name"`
}

type PropertyID string

func (pID PropertyID) String() string {
	return string(pID)
}
