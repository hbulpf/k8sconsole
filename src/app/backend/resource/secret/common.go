package secret

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"k8s.io/api/core/v1"
)

type SecretCell v1.Secret

func (self SecretCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
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

func toCells(std []v1.Secret) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = SecretCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []v1.Secret {
	std := make([]v1.Secret, len(cells))
	for i := range std {
		std[i] = v1.Secret(cells[i].(SecretCell))
	}
	return std
}
