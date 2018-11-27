package replicationcontroller

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type ReplicationControllerList struct {
	ListMeta api.ListMeta `json:"listMeta"`

	// Basic information about resources status on the list
	Status common.ResourceStatus `json:"status"`

	// Unordered list of Replication Controllers
	ReplicationControllers []ReplicationController `json:"replicationControllers"`

	// List of non-critical errors, that occurred during resource retrieval
	Errors []error `json:"errors"`
}

type ReplicationController struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta api.TypeMeta `json:"typeMeta"`

	// Aggregate information about pods belonging to this Replication Controller
	Pods common.PodInfo `json:"pods"`

	// Container images of the Replication Controller
	ContainerImages []string `json:"containerImages"`

	// Init Container images of the Replication Controller
	InitContainerImages []string `json:"initContainerImages"`
}

func GetReplicationControllerList(client kubernetes.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*ReplicationControllerList, error) {
	glog.Info("Getting list of all replication controllers in the cluster")

	channels := &common.ResourceChannels{
		ReplicationControllerList: common.GetReplicationControllerListChannel(client, nsQuery, 1),
		PodList:                   common.GetPodListChannel(client, nsQuery, 1),
		EventList:                 common.GetEventListChannel(client, nsQuery, 1),
	}

	return GetReplicationControllerListFromChannels(channels, dsQuery)
}

func GetReplicationControllerListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery) (
	*ReplicationControllerList, error) {

	rcList := <- channels.ReplicationControllerList.List
	err := <- channels.ReplicationControllerList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	podList := <- channels.PodList.List
	err = <- channels.PodList.Error
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	eventList := <- channels.EventList.List
	err = <- channels.EventList.Error
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	rcs := toReplicationControllerList(rcList.Items, dsQuery, podList.Items, eventList.Items, nonCriticalErrors)
	rcs.Status = getStatus(rcList, podList.Items, eventList.Items)
	return rcs, nil
}

func toReplicationControllerList(rcs []v1.ReplicationController, dsQuery *dataselect.DataSelectQuery,
	pods []v1.Pod, events []v1.Event, nonCriticalErrors []error) *ReplicationControllerList {

	rcList := &ReplicationControllerList{
		ReplicationControllers: make([]ReplicationController, 0),
		ListMeta: api.ListMeta{TotalItems: len(rcs)},
		Errors: nonCriticalErrors,
	}
	rcCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(rcs), dsQuery)
	rcs = fromCells(rcCells)
	rcList.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, rc := range rcs {
		matchingPods := common.FilterPodsByControllerRef(&rc, pods)
		podInfo := common.GetPodInfo(rc.Status.Replicas, rc.Spec.Replicas, matchingPods)
		podInfo.Warnings = event.GetPodsEventWarnings(events, matchingPods)

		replicationController := ToReplicationController(&rc, &podInfo)
		rcList.ReplicationControllers = append(rcList.ReplicationControllers, replicationController)
	}

	return rcList
}

func ToReplicationController(rc *v1.ReplicationController, podInfo *common.PodInfo) ReplicationController {
	return ReplicationController{
		ObjectMeta: api.NewObjectMeta(rc.ObjectMeta),
		TypeMeta: api.NewTypeMeta(api.ResourceKindReplicationController),
		Pods: *podInfo,
		ContainerImages: common.GetContainerImages(&rc.Spec.Template.Spec),
		InitContainerImages: common.GetInitContainerImages(&rc.Spec.Template.Spec),
	}
}