package notionapi

type FilterOperator string

const (
	FilterOperatorAND FilterOperator = "and"
	FilterOperatorOR  FilterOperator = "or"
)

type Condition string

const (
	ConditionEquals         Condition = "equals"
	ConditionDoesNotEqual   Condition = "does_not_equal"
	ConditionContains       Condition = "contains"
	ConditionDoesNotContain Condition = "does_not_contain"
	ConditionDoesStartsWith Condition = "starts_with"
	ConditionDoesEndsWith   Condition = "ends_with"
	ConditionDoesIsEmpty    Condition = "is_empty"
	ConditionGreaterThan    Condition = "greater_than"
	ConditionLessThan       Condition = "less_than"

	ConditionGreaterThanOrEqualTo Condition = "greater_than_or_equal_to"
	ConditionLessThanOrEqualTo    Condition = "greater_than_or_equal_to"

	ConditionBefore     Condition = "before"
	ConditionAfter      Condition = "after"
	ConditionOnOrBefore Condition = "on_or_before"
	ConditionOnOrAfter  Condition = "on_or_after"
	ConditionPastWeek   Condition = "past_week"
	ConditionPastMonth  Condition = "past_month"
	ConditionPastYear   Condition = "past_year"
	ConditionNextWeek   Condition = "next_week"
	ConditionNextMonth  Condition = "next_month"
	ConditionNextYear   Condition = "next_year"

	ConditionText     Condition = "text"
	ConditionCheckbox Condition = "checkbox"
	ConditionNumber   Condition = "number"
	ConditionDate     Condition = "date"
)

type FilterObject interface{}
