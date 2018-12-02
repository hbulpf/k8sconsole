package replicaset

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	ds "github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	hpa "github.com/wzt3309/k8sconsole/src/app/backend/resource/horizontalpodautoscaler"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/pod"
	resourceService "github.com/wzt3309/k8sconsole/src/app/backend/resource/service"
	apps "k8s.io/api/apps/v1beta2"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ReplicaSetDetail is a presentation layer view of Kubernetes Replica Set resource. This means
// it is Replica Set plus additional augmented data we can get from other sources
// (like services that target the same pods).
type ReplicaSetDetail struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`

	// Aggregate information about pods belonging to this Replica Set.
	PodInfo common.PodInfo `json:"podInfo"`

	// Detailed information about Pods belonging to this Replica Set.
	PodList pod.PodList `json:"podList"`

	// Detailed information about service related to Replica Set.
	ServiceList resourceService.ServiceList `json:"serviceList"`

	// Container images of the Replica Set.
	ContainerImages []string `json:"containerImages"`

	// Init Container images of the Replica Set.
	InitContainerImages []string `json:"initContainerImages"`

	// List of events related to this Replica Set.
	EventList common.EventList `json:"eventList"`

	// Selector of this replica set.
	Selector *metaV1.LabelSelector `json:"selector"`

	// List of Horizontal Pod Autoscalers targeting this Replica Set.
	HorizontalPodAutoscalerList hpa.HorizontalPodAutoscalerList `json:"horizontalPodAutoscalerList"`

	// List of non-critical errors, that occurred during resource retrieval.
	Errors []error `json:"errors"`
}

func GetReplicaSetDetail(client kubernetes.Interface, 	namespace, name string) (*ReplicaSetDetail, error) {
	glog.Info("Getting details of %s service in %s namespace", name, namespace)

	rs, err := client.AppsV1beta2().ReplicaSets(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	eventList, err := event.GetResourceEvents(client, ds.DefaultDataSelect, rs.Namespace, rs.Name)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	podList, err := GetReplicaSetPods(client, ds.DefaultDataSelect, name, namespace)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	podInfo, err := getReplicaSetPodInfo(client, rs)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	serviceList, err := GetReplicaSetServices(client, ds.DefaultDataSelect, namespace, name)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	hpas, err := hpa.GetHorizontalPodAutoscalerListForResource(client, namespace, "ReplicaSet", name)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	rsDetail := toReplicaSetDetail(rs, *eventList, *podList, *podInfo, *serviceList, *hpas, nonCriticalErrors)
	return &rsDetail, nil
}

func toReplicaSetDetail(replicaSet *apps.ReplicaSet, eventList common.EventList, podList pod.PodList,
	podInfo common.PodInfo, serviceList resourceService.ServiceList,
	hpas hpa.HorizontalPodAutoscalerList, nonCriticalErrors []error) ReplicaSetDetail {

	return ReplicaSetDetail{
		ObjectMeta:                  api.NewObjectMeta(replicaSet.ObjectMeta),
		TypeMeta:                    api.NewTypeMeta(api.ResourceKindReplicaSet),
		ContainerImages:             common.GetContainerImages(&replicaSet.Spec.Template.Spec),
		InitContainerImages:         common.GetInitContainerImages(&replicaSet.Spec.Template.Spec),
		Selector:                    replicaSet.Spec.Selector,
		PodInfo:                     podInfo,
		PodList:                     podList,
		ServiceList:                 serviceList,
		EventList:                   eventList,
		HorizontalPodAutoscalerList: hpas,
		Errors: nonCriticalErrors,
	}
}