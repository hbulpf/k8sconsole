package dataselect

import "sort"

type DataCell interface {
	GetProperty(PropertyName) ComparableValue
}

// ComparableValue hold any value that can be compared ot its own kind.
type ComparableValue interface {
	Compare(ComparableValue) int
	Contains(ComparableValue) bool
}

// DataSelector contains all the required data to perform data selection.
// Implements Sort interface
type DataSelector struct {
	GenericDataList	[]DataCell
	DataSelectQuery *DataSelectQuery
}

// Len returns the length of data inside DataSelector.
func (self DataSelector) Len() int {
	return len(self.GenericDataList)
}

// Swap swaps 2 indices inside DataSelector.
func (self DataSelector) Swap(i, j int) {
	self.GenericDataList[i], self.GenericDataList[j] = self.GenericDataList[j], self.GenericDataList[i]
}

// Less compares 2 indices, returns true if the element i should be sorted before j
func (self DataSelector) Less(i, j int) bool {
	for _, sortBy := range self.DataSelectQuery.SortQuery.SortByList {
		a := self.GenericDataList[i].GetProperty(sortBy.Property)
		b := self.GenericDataList[j].GetProperty(sortBy.Property)

		// ignore if property not found
		if a == nil || b == nil {
			break
		}

		cmp := a.Compare(b)
		if cmp == 0 {
			continue
		} else {
			return (cmp == -1 && sortBy.Asc) || (cmp == 1 && !sortBy.Asc)
		}
	}
	return false
}

func (self *DataSelector) Sort() *DataSelector {
	sort.Sort(*self)
	return self
}

// Filter the data inside
func (self *DataSelector) Filter() *DataSelector {
	filteredList := []DataCell{}

	for _, c := range self.GenericDataList {
		matches := true
		for _, filterBy := range self.DataSelectQuery.FilterQuery.FilterByList {
			v := c.GetProperty(filterBy.Property)
			// if cannot find property, continue
			if v == nil {
				matches = false
				continue
			}

			if !v.Contains(filterBy.Value) {
				matches = false
				continue
			}
		}

		if matches {
			filteredList = append(filteredList, c)
		}
	}

	self.GenericDataList = filteredList
	return self
}

// Paginates the data inside
func (self *DataSelector) Paginate() *DataSelector {
	pQuery := self.DataSelectQuery.PaginationQuery
	// return all items if provided settings do not meet requirements
	if !pQuery.IsValidPagination() {
		return self
	}

	dataList := self.GenericDataList
	startIndex, endIndex := pQuery.GetPaginationSettings(len(dataList))

	// return empty if required page does not exist
	if !pQuery.IsPageAvaliable(len(dataList), startIndex) {
		self.GenericDataList = []DataCell{}
		return self
	}

	self.GenericDataList = dataList[startIndex:endIndex]
	return self
}

// GenericDataSelect takes a list of DataCells and DtaSelectQuery and returns selected data.
func GenericDataSelect(dataList []DataCell, dsQuery *DataSelectQuery) []DataCell {
	SelectableData := DataSelector{
		GenericDataList: dataList,
		DataSelectQuery: dsQuery,
	}

	return SelectableData.Sort().Paginate().GenericDataList
}

// GenericDataSelectWithFilter takes a list of DataCells and DtaSelectQuery and returns selected data.
func GenericDataSelectWithFilter(dataList []DataCell, dsQuery *DataSelectQuery) ([]DataCell, int) {
	SelectableData := DataSelector{
		GenericDataList: dataList,
		DataSelectQuery: dsQuery,
	}

	filtered := SelectableData.Filter()
	filteredTotal := len(filtered.GenericDataList)
	processed := filtered.Sort().Paginate()
	return processed.GenericDataList, filteredTotal
}