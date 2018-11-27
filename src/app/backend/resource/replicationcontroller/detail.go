package replicationcontroller

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	ds "github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/pod"
	resourceService "github.com/wzt3309/k8sconsole/src/app/backend/resource/service"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ReplicationControllerDetail represents detailed information about a Replication Controller.
type ReplicationControllerDetail struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`

	// Label selector of the Replication Controller.
	LabelSelector map[string]string `json:"labelSelector"`

	// Container image list of the pod template specified by this Replication Controller.
	ContainerImages []string `json:"containerImages"`

	// Init Container image list of the pod template specified by this Replication Controller.
	InitContainerImages []string `json:"initContainerImages"`

	// Aggregate information about pods of this replication controller.
	PodInfo common.PodInfo `json:"podInfo"`

	// Detailed information about Pods belonging to this Replication Controller.
	PodList pod.PodList `json:"podList"`

	// Detailed information about service related to Replication Controller.
	ServiceList resourceService.ServiceList `json:"serviceList"`

	// List of events related to this Replication Controller.
	EventList common.EventList `json:"eventList"`

	// True when the data contains at least one pod with metrics information, false otherwise.
	HasMetrics bool `json:"hasMetrics"`

	// List of non-critical errors, that occurred during resource retrieval.
	Errors []error `json:"errors"`
}

// ReplicationControllerSpec contains information needed to update replication controller.
type ReplicationControllerSpec struct {
	// Replicas (pods) number in replicas set
	Replicas int32 `json:"replicas"`
}

// GetReplicationControllerDetail returns detailed information about the given replication controller
// in the given namespace.
func GetReplicationControllerDetail(client kubernetes.Interface, namespace, name string) (
	*ReplicationControllerDetail, error) {
	glog.Info("Getting details of %s replication controller in %s namespace", name, namespace)

	replicationController, err := client.CoreV1().ReplicationControllers(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	podInfo, err := getReplicationControllerPodInfo(client, replicationController, namespace)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	podList, err := GetReplicationControllerPods(client, ds.DefaultDataSelect, name, namespace)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	eventList, err := event.GetResourceEvents(client, ds.DefaultDataSelect, namespace, name)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	serviceList, err := GetReplicationControllerServices(client, ds.DefaultDataSelect, namespace, name)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}


	replicationControllerDetail := toReplicationControllerDetail(replicationController, *podInfo,
		*podList, *eventList, *serviceList, nonCriticalErrors)
	return &replicationControllerDetail, nil
}

func UpdateReplicasCount(client kubernetes.Interface, namespace, name string, spec *ReplicationControllerSpec) error {
	glog.Infof("Updating replicas count to %d for %s replication controller in %s namespace",
		spec.Replicas, name, namespace)

	rc, err := client.CoreV1().ReplicationControllers(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		return err
	}

	rc.Spec.Replicas = &spec.Replicas

	_, err = client.CoreV1().ReplicationControllers(namespace).Update(rc)
	if err != nil {
		return nil
	}

	glog.Infof("Successfully update replicas count to %d for %s replication controller from %s namespace",
		spec.Replicas, name, namespace)

	return nil
}

func toReplicationControllerDetail(replicationController *v1.ReplicationController, podInfo common.PodInfo,
	podList pod.PodList, eventList common.EventList, serviceList resourceService.ServiceList,
	nonCriticalErrors []error) ReplicationControllerDetail {

	return ReplicationControllerDetail{
		ObjectMeta:                  api.NewObjectMeta(replicationController.ObjectMeta),
		TypeMeta:                    api.NewTypeMeta(api.ResourceKindReplicationController),
		LabelSelector:               replicationController.Spec.Selector,
		PodInfo:                     podInfo,
		PodList:                     podList,
		EventList:                   eventList,
		ServiceList:                 serviceList,
		ContainerImages:             common.GetContainerImages(&replicationController.Spec.Template.Spec),
		InitContainerImages:         common.GetInitContainerImages(&replicationController.Spec.Template.Spec),
		Errors:                      nonCriticalErrors,
	}
}