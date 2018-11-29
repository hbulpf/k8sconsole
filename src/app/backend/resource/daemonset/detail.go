package daemonset

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	ds "github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/pod"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/service"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// DaemonSeDetail represents detailed information about a Daemon Set.
type DaemonSetDetail struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`

	// Label selector of the Daemon Set.
	LabelSelector *v1.LabelSelector `json:"labelSelector,omitempty"`

	// Container image list of the pod template specified by this Daemon Set.
	ContainerImages []string `json:"containerImages"`

	// Init Container image list of the pod template specified by this Daemon Set.
	InitContainerImages []string `json:"initContainerImages"`

	// Aggregate information about pods of this daemon set.
	PodInfo common.PodInfo `json:"podInfo"`

	// Detailed information about Pods belonging to this Daemon Set.
	PodList pod.PodList `json:"podList"`

	// Detailed information about service related to Daemon Set.
	ServiceList service.ServiceList `json:"serviceList"`

	// True when the data contains at least one pod with metrics information, false otherwise.
	HasMetrics bool `json:"hasMetrics"`

	// List of events related to this daemon set
	EventList common.EventList `json:"eventList"`

	// List of non-critical errors, that occurred during resource retrieval.
	Errors []error `json:"errors"`
}

// Returns detailed information about the given daemon set in the given namespace.
func GetDaemonSetDetail(client kubernetes.Interface, namespace, name string) (*DaemonSetDetail, error) {

	glog.Infof("Getting details of %s daemon set in %s namespace", name, namespace)
	daemonSet, err := client.AppsV1beta2().DaemonSets(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	podList, err := GetDaemonSetPods(client, ds.DefaultDataSelect, name, namespace)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	podInfo, err := getDaemonSetPodInfo(client, daemonSet)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	serviceList, err := GetDaemonSetServices(client, ds.DefaultDataSelect, namespace, name)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	eventList, err := event.GetResourceEvents(client, ds.DefaultDataSelect, daemonSet.Namespace, daemonSet.Name)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	daemonSetDetail := &DaemonSetDetail{
		ObjectMeta:    api.NewObjectMeta(daemonSet.ObjectMeta),
		TypeMeta:      api.NewTypeMeta(api.ResourceKindDaemonSet),
		LabelSelector: daemonSet.Spec.Selector,
		PodInfo:       *podInfo,
		PodList:       *podList,
		ServiceList:   *serviceList,
		EventList:     *eventList,
		Errors:        nonCriticalErrors,
	}

	for _, container := range daemonSet.Spec.Template.Spec.Containers {
		daemonSetDetail.ContainerImages = append(daemonSetDetail.ContainerImages, container.Image)
	}

	for _, initContainer := range daemonSet.Spec.Template.Spec.InitContainers {
		daemonSetDetail.InitContainerImages = append(daemonSetDetail.InitContainerImages, initContainer.Image)
	}

	return daemonSetDetail, nil
}
