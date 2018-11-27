package replicationcontroller

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// Transforms simple selector map to labels.Selector object that can be used when querying for object
func toLabelSelector(selector map[string]string) (labels.Selector, error) {
	labelSelector, err := metaV1.LabelSelectorAsSelector(&metaV1.LabelSelector{MatchLabels: selector})

	if err != nil {
		return nil, err
	}

	return labelSelector, nil
}

// Services are matched by replication controllers' label selector. They are deleted if give
// label selector is targeting only 1 replication controller
func getServicesForDeletion(client kubernetes.Interface, labelSelector labels.Selector,
	namespace string) ([]v1.Service, error) {
	replicationControllers, err := client.CoreV1().ReplicationControllers(namespace).List(metaV1.ListOptions{
		LabelSelector: labelSelector.String(),
		FieldSelector: fields.Everything().String(),
	})
	if err != nil {
		return nil, err
	}

	// If label selector is targeting only 1 replication controller
	// then we can delete services targeted by this label selector. They are deleted if give
	// label selector is targeting only 1 replication controller
	if len(replicationControllers.Items) != 1 {
		return []v1.Service{}, nil
	}

	services, err := client.CoreV1().Services(namespace).List(metaV1.ListOptions{
		LabelSelector: labelSelector.String(),
		FieldSelector: fields.Everything().String(),
	})
	if err != nil {
		return nil, err
	}

	return services.Items, nil
}

func getStatus(list *v1.ReplicationControllerList, pods []v1.Pod, events []v1.Event) common.ResourceStatus {
	info := common.ResourceStatus{}
	if list == nil {
		return info
	}

	for _, rc := range list.Items {
		matchingPods := common.FilterPodsByControllerRef(&rc, pods)
		podInfo := common.GetPodInfo(rc.Status.Replicas, rc.Spec.Replicas, matchingPods)
		warnings := event.GetPodsEventWarnings(events, matchingPods)

		if len(warnings) > 0 {
			info.Failed++
		} else if podInfo.Pending > 0 {
			info.Pending++
		} else {
			info.Running++
		}
	}
	return info
}

type ReplicationControllerCell v1.ReplicationController

func (self ReplicationControllerCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Name)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Namespace)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(self.ObjectMeta.CreationTimestamp.Time)
	default:
		return nil
	}
}

func toCells(std []v1.ReplicationController) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = ReplicationControllerCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []v1.ReplicationController {
	std := make([]v1.ReplicationController, len(cells))
	for i := range cells {
		std[i] = v1.ReplicationController(cells[i].(ReplicationControllerCell))
	}
	return std
}

