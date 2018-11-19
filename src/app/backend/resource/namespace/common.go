package namespace

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// NamespaceSpec is a specification of namespace to create
type NamespaceSpec struct {
	// Name of the namespace
	Name string `json:"name"`
}

// CreateNamespace creates namespace based on given specification.
func CreateNamespace(spec *NamespaceSpec, client kubernetes.Interface) error {
	glog.Infof("Creating namespace %s", spec.Name)

	namespace := &v1.Namespace{
		ObjectMeta: metaV1.ObjectMeta{
			Name: spec.Name,
		},
	}

	_, err := client.CoreV1().Namespaces().Create(namespace)
	return err
}

type NamespaceCell v1.Namespace

func (self NamespaceCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
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

func toCells(std []v1.Namespace) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = NamespaceCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []v1.Namespace {
	std := make([]v1.Namespace, len(cells))
	for i := range std {
		std[i] = v1.Namespace(cells[i].(NamespaceCell))
	}
	return std
}
