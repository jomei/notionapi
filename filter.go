package notionapi

type FilterOperator string

type CompoundFilter map[FilterOperator][]PropertyFilter

type Condition string

type PropertyFilter struct {
	Property    string                      `json:"property"`
	Text        *TextFilterCondition        `json:"text,omitempty"`
	Number      *NumberFilterCondition      `json:"number,omitempty"`
	Checkbox    *CheckboxFilterCondition    `json:"checkbox,omitempty"`
	Select      *SelectFilterCondition      `json:"select,omitempty"`
	MultiSelect *MultiSelectFilterCondition `json:"multi_select,omitempty"`
	Date        *DateFilterCondition        `json:"date,omitempty"`
	People      *PeopleFilterCondition      `json:"people,omitempty"`
	Files       *FilesFilterCondition       `json:"files,omitempty"`
	Relation    *RelationFilterCondition    `json:"relation,omitempty"`
	Formula     *FormulaFilterCondition     `json:"formula,omitempty"`
}

type TextFilterCondition struct {
	Equals         string `json:"equals,omitempty"`
	DoesNotEqual   string `json:"does_not_equal,omitempty"`
	Contains       string `json:"contains,omitempty"`
	DoesNotContain string `json:"does_not_contain,omitempty"`
	StartsWith     string `json:"starts_with,omitempty"`
	EndsWith       string `json:"ends_with,omitempty"`
	IsEmpty        bool   `json:"is_empty,omitempty"`
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`
}

type NumberFilterCondition struct {
	Equals               float64 `json:"equals,omitempty"`
	DoesNotEqual         float64 `json:"does_not_equal,omitempty"`
	GreaterThan          float64 `json:"greater_than,omitempty"`
	LessThan             float64 `json:"less_than,omitempty"`
	GreaterThanOrEqualTo float64 `json:"greater_than_or_equal_to"`
	LessThanOrEqualTo    float64 `json:"less_than_or_equal_to"`
	IsEmpty              bool    `json:"is_empty,omitempty"`
	IsNotEmpty           bool    `json:"is_not_empty,omitempty"`
}

type CheckboxFilterCondition struct {
	Equals       bool `json:"equals,omitempty"`
	DoesNotEqual bool `json:"does_not_equal,omitempty"`
}

type SelectFilterCondition struct {
	Equals       string `json:"equals,omitempty"`
	DoesNotEqual string `json:"does_not_equal,omitempty"`
	IsEmpty      bool   `json:"is_empty,omitempty"`
	IsNotEmpty   bool   `json:"is_not_empty,omitempty"`
}

type MultiSelectFilterCondition struct {
	Contains       string `json:"contains,omitempty"`
	DoesNotContain string `json:"does_not_contain,omitempty"`
	IsEmpty        bool   `json:"is_empty,omitempty"`
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`
}

type DateFilterCondition struct {
	Equals     *Date     `json:"equals,omitempty"`
	Before     *Date     `json:"before,omitempty"`
	After      *Date     `json:"after,omitempty"`
	OnOrBefore *Date     `json:"on_or_before,omitempty"`
	OnOrAfter  *Date     `json:"on_or_after,omitempty"`
	PastWeek   *struct{} `json:"past_week,omitempty"`
	PastMonth  *struct{} `json:"past_month,omitempty"`
	PastYear   *struct{} `json:"past_year,omitempty"`
	NextWeek   *struct{} `json:"next_week,omitempty"`
	NextMonth  *struct{} `json:"next_month,omitempty"`
	NextYear   *struct{} `json:"next_year,omitempty"`
	IsEmpty    bool      `json:"is_empty,omitempty"`
	IsNotEmpty bool      `json:"is_not_empty,omitempty"`
}

type PeopleFilterCondition struct {
	Contains       string `json:"contains,omitempty"`
	DoesNotContain string `json:"does_not_contain,omitempty"`
	IsEmpty        bool   `json:"is_empty,omitempty"`
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`
}

type FilesFilterCondition struct {
	IsEmpty    bool `json:"is_empty,omitempty"`
	IsNotEmpty bool `json:"is_not_empty,omitempty"`
}

type RelationFilterCondition struct {
	Contains       string `json:"contains,omitempty"`
	DoesNotContain string `json:"does_not_contain,omitempty"`
	IsEmpty        bool   `json:"is_empty,omitempty"`
	IsNotEmpty     bool   `json:"is_not_empty,omitempty"`
}

type FormulaFilterCondition struct {
	Text     *TextFilterCondition     `json:"text,omitempty"`
	Checkbox *CheckboxFilterCondition `json:"checkbox,omitempty"`
	Number   *NumberFilterCondition   `json:"number,omitempty"`
	Date     *DateFilterCondition     `json:"date,omitempty"`
}
