package notionapi

import "time"

type FilterOperator string

type Filter interface{}

type CompoundFilter map[FilterOperator]Filter

type Condition string

type PropertyFilter struct {
	Property    PropertyType                    `json:"property"`
	Text        map[Condition]string            `json:"text,omitempty"`
	Number      map[Condition]float64           `json:"number,omitempty"`
	Checkbox    map[Condition]bool              `json:"checkbox,omitempty"`
	Select      map[Condition]interface{}       `json:"select,omitempty"`
	MultiSelect map[Condition]interface{}       `json:"multi_select,omitempty"`
	Date        map[Condition]time.Time         `json:"date,omitempty"`
	People      map[Condition]interface{}       `json:"people,omitempty"`
	Files       map[Condition]bool              `json:"files,omitempty"`
	Relation    map[Condition]interface{}       `json:"relation,omitempty"`
	Formula     map[PropertyType]PropertyFilter `json:"formula,omitempty"`
}
