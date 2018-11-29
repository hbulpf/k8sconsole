package daemonset

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	apps "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// DaemonSetList contains a list of Daemon Sets in the cluster.
type DaemonSetList struct {
	ListMeta          api.ListMeta       `json:"listMeta"`
	DaemonSets        []DaemonSet        `json:"daemonSets"`

	// Basic information about resources status on the list.
	Status common.ResourceStatus `json:"status"`

	// List of non-critical errors, that occurred during resource retrieval.
	Errors []error `json:"errors"`
}

// DaemonSet plus zero or more Kubernetes services that target the Daemon Set.
type DaemonSet struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`

	// Aggregate information about pods belonging to this Daemon Set.
	Pods common.PodInfo `json:"pods"`

	// Container images of the Daemon Set.
	ContainerImages []string `json:"containerImages"`

	// InitContainer images of the Daemon Set.
	InitContainerImages []string `json:"initContainerImages"`
}

// GetDaemonSetList returns a list of all Daemon Set in the cluster.
func GetDaemonSetList(client kubernetes.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*DaemonSetList, error) {
	channels := &common.ResourceChannels{
		DaemonSetList: common.GetDaemonSetListChannel(client, nsQuery, 1),
		ServiceList:   common.GetServiceListChannel(client, nsQuery, 1),
		PodList:       common.GetPodListChannel(client, nsQuery, 1),
		EventList:     common.GetEventListChannel(client, nsQuery, 1),
	}

	return GetDaemonSetListFromChannels(channels, dsQuery)
}

// GetDaemonSetListFromChannels returns a list of all Daemon Set in the cluster
// reading required resource list once from the channels.
func GetDaemonSetListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery) (
	*DaemonSetList, error) {

	daemonSets := <-channels.DaemonSetList.List
	err := <-channels.DaemonSetList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	pods := <-channels.PodList.List
	err = <-channels.PodList.Error
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	events := <-channels.EventList.List
	err = <-channels.EventList.Error
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	dsList := toDaemonSetList(daemonSets.Items, pods.Items, events.Items, nonCriticalErrors, dsQuery)
	dsList.Status = getStatus(daemonSets, pods.Items, events.Items)
	return dsList, nil
}

func toDaemonSetList(daemonSets []apps.DaemonSet, pods []v1.Pod, events []v1.Event, nonCriticalErrors []error,
	dsQuery *dataselect.DataSelectQuery) *DaemonSetList {

	daemonSetList := &DaemonSetList{
		DaemonSets: make([]DaemonSet, 0),
		ListMeta:   api.ListMeta{TotalItems: len(daemonSets)},
		Errors:     nonCriticalErrors,
	}

	dsCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(daemonSets), dsQuery)
	daemonSets = fromCells(dsCells)
	daemonSetList.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for i, daemonSet := range daemonSets {
		matchingPods := common.FilterPodsByControllerRef(&daemonSet, pods)
		podInfo := common.GetPodInfo(daemonSet.Status.CurrentNumberScheduled,
			&daemonSets[i].Status.DesiredNumberScheduled, matchingPods)
		podInfo.Warnings = event.GetPodsEventWarnings(events, matchingPods)

		daemonSetList.DaemonSets = append(daemonSetList.DaemonSets, DaemonSet{
			ObjectMeta:          api.NewObjectMeta(daemonSet.ObjectMeta),
			TypeMeta:            api.NewTypeMeta(api.ResourceKindDaemonSet),
			Pods:                podInfo,
			ContainerImages:     common.GetContainerImages(&daemonSet.Spec.Template.Spec),
			InitContainerImages: common.GetInitContainerImages(&daemonSet.Spec.Template.Spec),
		})
	}

	return daemonSetList
}