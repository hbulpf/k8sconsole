package configmap

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"k8s.io/api/core/v1"
)

type ConfigMapCell v1.ConfigMap

func (self ConfigMapCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
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

func toCells(std []v1.ConfigMap) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = ConfigMapCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []v1.ConfigMap {
	std := make([]v1.ConfigMap, len(cells))
	for i := range std {
		std[i] = v1.ConfigMap(cells[i].(ConfigMapCell))
	}
	return std
}