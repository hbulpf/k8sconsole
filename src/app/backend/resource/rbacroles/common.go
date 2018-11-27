package rbacroles

import "github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"

type RoleCell RbacRole

func (self RoleCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Name)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(self.ObjectMeta.CreationTimestamp.Time)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Namespace)
	default:
		// if name is not supported then just return a constant dummy value, sort will have no effect.
		return nil
	}
}

func toCells(std []RbacRole) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = RoleCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []RbacRole {
	std := make([]RbacRole, len(cells))
	for i := range std {
		std[i] = RbacRole(cells[i].(RoleCell))
	}
	return std
}