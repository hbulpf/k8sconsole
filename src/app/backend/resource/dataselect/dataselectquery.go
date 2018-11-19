package dataselect

type DataSelectQuery struct {
	PaginationQuery 	*PaginationQuery
	SortQuery					*SortQuery
	FilterQuery				*FilterQuery
}

// NoDataSelect is an option for no data select (same data will be returned).
var NoDataSelect = NewDataSelectQuery(NoPagination, NoSort, NoFilter)

// DefaultDataSelect downloads first 10 items from page 1 with no sort.
var DefaultDataSelect = NewDataSelectQuery(DefaultPagination, NoSort, NoFilter)

//  NewDataSelectQuery creates DataSelectQuery object from data select queries
func NewDataSelectQuery(paginationQuery *PaginationQuery, sortQuery *SortQuery, filterQuery *FilterQuery) *DataSelectQuery {
	return &DataSelectQuery{
		PaginationQuery: paginationQuery,
		SortQuery: sortQuery,
		FilterQuery: filterQuery,
	}
}
