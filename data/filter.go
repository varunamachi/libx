package data

//FilterType - Type of filter item
type FilterType string

//Prop - filter for a value
const Prop FilterType = "Prop"

//Array - filter for an array
const Array FilterType = "arrayArray"

//DateRange - filter for date range
const DtRange FilterType = "DateRange"

//NumRange - filter for real number range
const NumRange FilterType = "NumRange"

//Boolean - filter for boolean field
const Boolean FilterType = "Boolean"

//Search - filter for search text field
const Search FilterType = "Search"

//Constant - constant filter value
const Constant FilterType = "Constant"

//Static - constant filter value
const Static FilterType = "Static"

//MatchStrategy - strategy to match multiple fields passed as part of the
//filters
type MatchStrategy string

//MatchAll - match all provided values while executing filter
const MatchAll MatchStrategy = "All"

//MatchOne - match atleast one of the  provided values while executing filter
const MatchOne MatchStrategy = "One"

//MatchNone - match values that are not part of the provided list while
//executing filter
const MatchNone MatchStrategy = "None"

//FilterSpec - filter specification
type FilterSpec struct {
	Field string     `json:"field" db:"field" bson:"field"`
	Name  string     `json:"name" db:"name" bson:"name"`
	Type  FilterType `json:"type" db:"type" bson:"type"`
}

//Matcher - matches the given fields. If multiple fileds are given the; the
//joining condition is decided by the MatchStrategy given
type Matcher struct {
	Strategy MatchStrategy `json:"strategy" db:"strategy" bson:"strategy"`
	Fields   []interface{} `json:"fields" db:"fields" bson:"fields"`
}

//PropMatcher - matches props
type PropMatcher []interface{}

//FilterSpecList - alias for array of filter specs
type FilterSpecList []*FilterSpec

//FilterVal - values for filter along with the count
type FilterVal struct {
	Name  string `json:"name" db:"name"`
	Count int    `json:"count" db:"count"`
}

//Filter - generic filter used to filter data in any mongodb collection
type Filter struct {
	Props    map[string]Matcher     `json:"props" db:"props"`
	Bools    map[string]interface{} `json:"bools" db:"bools"`
	Dates    map[string]DateRange   `json:"dates" db:"dates"`
	Lists    map[string]Matcher     `json:"lists" db:"lists"`
	Searches map[string]Matcher     `json:"searches" db:"searches"`
}
