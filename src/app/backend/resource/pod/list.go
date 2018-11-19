package pod

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sClient "k8s.io/client-go/kubernetes"
)

type PodList struct {
	ListMeta 		api.ListMeta 						`json:"listMeta"`

	// Basic information about resources status on th list
	Status 			common.ResourceStatus		`json:"status"`

	Pods 				[]Pod										`json:"pods"`

	// List of non-critical errors, that occurred during resource retrieval
	Errors 			[]error									`json:"errors"`
}

type PodStatus struct {
	Status 						string 								`json:"status"`
	PodPhase 					v1.PodPhase 					`json:"podPhase"`
	ContainerStates 	[]v1.ContainerState		`json:"containerStates"`
}

// Pod is a view of kubernetes Pod resource.
// It is Pod plus additional augmented data
type Pod struct {
	ObjectMeta 		api.ObjectMeta 	`json:"objectMeta"`
	TypeMeta   		api.TypeMeta   	`json:"typeMeta"`

	// more info on pod status
	PodStatus			PodStatus				`json:"podStatus"`

	// Count of containers restarts
	RestartCount	int32 					`json:"restartCount"`

	// Name of the node this pod runs on
	NodeName 			string					`json:"nodeName"`
}

var EmptyPodList = &PodList{
	Pods: make([]Pod, 0),
	Errors: make([]error, 0),
	ListMeta: api.ListMeta{
		TotalItems: 0,
	},
}

// GetPodList returns a list of all Pods in cluster
func GetPodList(client k8sClient.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*PodList, error) {
	glog.Info("Getting list of all pods in cluster")

	channels := &common.ResourceChannels{
		PodList: common.GetPodListChannelWithOptions(client, nsQuery, metaV1.ListOptions{}, 1),
		EventList: common.GetEventListChannel(client, nsQuery, 1),
	}

	return GetPodListFromChannels(channels, dsQuery)
}

// GetPodListFromChannels returns a list of all Pods in the cluster
func GetPodListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery) (*PodList, error) {

	pods := <- channels.PodList.List
	err := <- channels.PodList.Error
	nonCriticalErrors, criticalErrors := errors.HandleError(err)
	if criticalErrors != nil {
		return nil, criticalErrors
	}

	eventList := <- channels.EventList.List
	err = <- channels.EventList.Error
	nonCriticalErrors, criticalErrors = errors.AppendError(err, nonCriticalErrors)
	if criticalErrors != nil {
		return nil, criticalErrors
	}

	podList := ToPodList(pods.Items, eventList.Items, nonCriticalErrors, dsQuery)
	podList.Status = getStatus(pods, eventList.Items)
	return &podList, nil
}

func ToPodList(pods []v1.Pod, events []v1.Event, nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery) PodList {
	podList := PodList{
		Pods: make([]Pod, 0),
		Errors: nonCriticalErrors,
	}

	// filter and sort pods
	podCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(pods), dsQuery)
	pods = fromCells(podCells)
	podList.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, pod := range pods {
		warnings := event.GetPodsEventWarnings(events, []v1.Pod{pod})
		podDetail := toPod(&pod, warnings)
		podList.Pods = append(podList.Pods, podDetail)
	}

	return podList
}

func toPod(pod *v1.Pod, warnings []common.Event) Pod {
	podDetail := Pod{
		ObjectMeta: 	api.NewObjectMeta(pod.ObjectMeta),
		TypeMeta: 		api.NewTypeMeta(api.ResourceKindPod),
		PodStatus: 		getPodStatus(*pod, warnings),
		RestartCount:	getRestartCount(*pod),
		NodeName: 		pod.Spec.NodeName,
	}

	return podDetail
}