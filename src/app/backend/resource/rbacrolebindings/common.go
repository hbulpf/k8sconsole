package rbacrolebindings

import "github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"

type RoleBindingCell RbacRoleBinding

func (self RoleBindingCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
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

func toCells(std []RbacRoleBinding) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = RoleBindingCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []RbacRoleBinding {
	std := make([]RbacRoleBinding, len(cells))
	for i := range cells {
		std[i] = RbacRoleBinding(cells[i].(RoleBindingCell))
	}
	return std
}

