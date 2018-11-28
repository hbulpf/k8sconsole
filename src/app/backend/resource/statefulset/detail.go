package statefulset

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	ds "github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/pod"
	apps "k8s.io/api/apps/v1beta2"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// StatefulSetDetail is a presentation layer view of Kubernetes Stateful Set resource. This means it is Stateful
// Set plus additional augmented data we can get from other sources (like services that target the same pods).
type StatefulSetDetail struct {
	ObjectMeta          api.ObjectMeta   `json:"objectMeta"`
	TypeMeta            api.TypeMeta     `json:"typeMeta"`
	PodInfo             common.PodInfo   `json:"podInfo"`
	PodList             pod.PodList      `json:"podList"`
	ContainerImages     []string         `json:"containerImages"`
	InitContainerImages []string         `json:"initContainerImages"`
	EventList           common.EventList `json:"eventList"`

	// List of non-critical errors, that occurred during resource retrieval.
	Errors []error `json:"errors"`
}

// GetStatefulSetDetail gets Stateful Set details.
func GetStatefulSetDetail(client kubernetes.Interface, namespace, name string) (*StatefulSetDetail, error) {
	glog.Infof("Getting details of %s statefulset in %s namespace", name, namespace)

	ss, err := client.AppsV1beta2().StatefulSets(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	podList, err := GetStatefulSetPods(client, ds.DefaultDataSelect, name, namespace)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	podInfo, err := getStatefulSetPodInfo(client, ss)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	events, err := event.GetResourceEvents(client, ds.DefaultDataSelect, ss.Namespace, ss.Name)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	ssDetail := getStatefulSetDetail(ss, *events, *podList, *podInfo, nonCriticalErrors)
	return &ssDetail, nil
}

func getStatefulSetDetail(statefulSet *apps.StatefulSet, eventList common.EventList, podList pod.PodList,
	podInfo common.PodInfo, nonCriticalErrors []error) StatefulSetDetail {
	return StatefulSetDetail{
		ObjectMeta:          api.NewObjectMeta(statefulSet.ObjectMeta),
		TypeMeta:            api.NewTypeMeta(api.ResourceKindStatefulSet),
		ContainerImages:     common.GetContainerImages(&statefulSet.Spec.Template.Spec),
		InitContainerImages: common.GetInitContainerImages(&statefulSet.Spec.Template.Spec),
		PodInfo:             podInfo,
		PodList:             podList,
		EventList:           eventList,
		Errors:              nonCriticalErrors,
	}
}