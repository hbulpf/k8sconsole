package ingress

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	extensions "k8s.io/api/extensions/v1beta1"
)

type IngressCell extensions.Ingress

func (self IngressCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Name)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(self.ObjectMeta.CreationTimestamp.Time)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Namespace)
	default:
		return nil
	}
}

func toCells(std []extensions.Ingress) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = IngressCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []extensions.Ingress {
	std := make([]extensions.Ingress, len(cells))
	for i := range std {
		std[i] = extensions.Ingress(cells[i].(IngressCell))
	}
	return std
}
