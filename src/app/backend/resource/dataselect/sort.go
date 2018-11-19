package dataselect

// SortQuery holds options for sort data
type SortQuery struct {
	SortByList	[]SortBy
}

// SortBy holds the name of property that should be sorted and whether order should be asc or desc.
type SortBy struct {
	Property   		PropertyName
	Asc						bool
}

var NoSort = &SortQuery{
	SortByList: []SortBy{},
}

// NewSortQuery takes raw sort options list and returns SortQuery object.
// e.g. ["a", "param1", "d", "param2"] meas sort by param1 asc - than the result sort by param2 desc
func NewSortQuery(sortByListRaw []string) *SortQuery {
	if sortByListRaw == nil || len(sortByListRaw) % 2 == 1 {
		return NoSort
	}

	sortByList := []SortBy{}
	for i := 0; i+1 < len(sortByListRaw); i+=2 {
		var asc bool
		orderOption := sortByListRaw[i]
		switch orderOption {
		case "a":
			asc = true
		case "d":
			asc = false
		default:
			return NoSort
		}

		propertyName := sortByListRaw[i+1]
		sortBy := SortBy{
			Property: PropertyName(propertyName),
			Asc: 			asc,
		}

		sortByList = append(sortByList, sortBy)
	}

	return &SortQuery{
		SortByList: sortByList,
	}
}
