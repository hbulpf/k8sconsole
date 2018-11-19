package dataselect

// By default backend pagination will not be applied.
var NoPagination = NewPaginationQuery(-1, -1)

// No items will be returned
var EmptyPagination = NewPaginationQuery(0, 0)

// Returns 10 items from page 1
var DefaultPagination = NewPaginationQuery(10, 0)

// PaginationQuery structure represents pagination settings
type PaginationQuery struct {
	// How many items per page should be returned
	ItemsPerPage	int
	// Number of page that should be returned
	Page int
}

func NewPaginationQuery(itemsPerPage, page int) *PaginationQuery {
	return &PaginationQuery{itemsPerPage, page}
}

func (p *PaginationQuery) IsValidPagination() bool {
	return p.ItemsPerPage >= 0 && p.Page >= 0
}

// IsPageAvailable returns true if at least one element can be placed on page.
func (p *PaginationQuery) IsPageAvaliable(itemsCount, startingIndex int) bool {
	return itemsCount > startingIndex && p.ItemsPerPage > 0
}

// Item index in page >= startIndex and < endIndex
func (p *PaginationQuery) GetPaginationSettings(itemsCount int) (startIndex int, endIndex int) {
	startIndex = p.ItemsPerPage * p.Page
	endIndex = startIndex + p.ItemsPerPage

	if endIndex > itemsCount {
		endIndex = itemsCount
	}

	return startIndex, endIndex
}
