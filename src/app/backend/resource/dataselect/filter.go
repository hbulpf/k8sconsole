package dataselect

type FilterQuery struct {
	FilterByList []FilterBy
}

type FilterBy struct {
	Property 	PropertyName
	Value 		ComparableValue
}

var NoFilter = &FilterQuery{
	FilterByList: []FilterBy{},
}

// NewFilterQuery takes raw filter options list and returns FilterQuery.
// e.g.
//   ["param1", "val1", "param2", "val2"] - means that the data should be
//   filtered by param1=val1, param2=va2
func NewFilterQuery(filterByListRaw []string) *FilterQuery {
	if filterByListRaw == nil || len(filterByListRaw)%2 == 1 {
		return NoFilter
	}
	filterByList := []FilterBy{}
	for i := 0; i+1 < len(filterByListRaw); i+=2 {
		propertyName := filterByListRaw[i]
		propertyValue := filterByListRaw[i+1]
		filterBy := FilterBy{
			Property: PropertyName(propertyName),
			Value: StdComparableString(propertyValue),
		}
		filterByList = append(filterByList, filterBy)
	}
	return &FilterQuery{
		FilterByList: filterByList,
	}
}

