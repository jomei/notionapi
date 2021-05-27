package notionapi

const (
	ObjectTypeDatabase ObjectType = "database"
	ObjectTypeBlock    ObjectType = "block"
	ObjectTypePage     ObjectType = "page"
)

const (
	PropertyTypeTitle          PropertyType = "title"
	PropertyTypeRichText       PropertyType = "rich_text"
	PropertyTypeSelect         PropertyType = "select"
	PropertyTypeNumber         PropertyType = "number"
	PropertyTypeCheckbox       PropertyType = "checkbox"
	PropertyTypeEmail          PropertyType = "email"
	PropertyTypeURL            PropertyType = "url"
	PropertyTypeFile           PropertyType = "file"
	PropertyTypePhoneNumber    PropertyType = "phone_number"
	PropertyTypeFormula        PropertyType = "formula"
	PropertyTypeDate           PropertyType = "date"
	PropertyTypeRelation       PropertyType = "relation"
	PropertyTypeRollup         PropertyType = "rollup"
	PropertyTypeMultiSelect    PropertyType = "multi_select"
	PropertyTypePeople         PropertyType = "people"
	PropertyTypeCreatedTime    PropertyType = "created_time"
	PropertyTypeCreatedBy      PropertyType = "created_by"
	PropertyTypeLastEditedTime PropertyType = "last_edited_time"
	PropertyTypeLastEditedBy   PropertyType = "last_edited_by"
)

const (
	FormatNumber           FormatType = "number"
	FormatNumberWithCommas FormatType = "number_with_commas"
	FormatPercent          FormatType = "percent"
	FormatDollar           FormatType = "dollar"
	FormatEuro             FormatType = "euro"
	FormatPound            FormatType = "pound"
	FormatYen              FormatType = "yen"
	FormatRuble            FormatType = "ruble"
	FormatRupee            FormatType = "rupee"
	FormatYuan             FormatType = "yuan"
)

const (
	ColorDefault Color = "default"
	ColorGray    Color = "gray"
	ColorBrown   Color = "brown"
	ColorOrange  Color = "orange"
	ColorYellow  Color = "yellow"
	ColorGreen   Color = "green"
	ColorBlue    Color = "blue"
	ColorPurple  Color = "purple"
	ColorPink    Color = "pink"
	ColorRed     Color = "red"
)

const (
	FunctionCountAll          FunctionType = "count_all"
	FunctionCountValues       FunctionType = "count_values"
	FunctionCountUniqueValues FunctionType = "count_unique_values"
	FunctionCountEmpty        FunctionType = "count_empty"
	FunctionCountNotEmpty     FunctionType = "count_not_empty"
	FunctionPercentEmpty      FunctionType = "percent_empty"
	FunctionPercentNotEmpty   FunctionType = "percent_not_empty"
	FunctionSum               FunctionType = "sum"
	FunctionAverage           FunctionType = "average"
	FunctionMedian            FunctionType = "median"
	FunctionMin               FunctionType = "min"
	FunctionMax               FunctionType = "max"
	FunctionRange             FunctionType = "range"
)
