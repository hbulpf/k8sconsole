package deployment

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	apps "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentList struct {
	ListMeta api.ListMeta `json:"listMeta"`

	// Basic information about resource status on the list
	Status common.ResourceStatus `json:"status"`

	// Unordered list of Deployments
	Deployments []Deployment `json:"deployments"`

	// List of non-critical errors, that occurred during resource retrieval
	Errors []error `json:"errors"`
}

// Deployment is a presentation layer view of Kubernetes Deployment resource.
type Deployment struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta api.TypeMeta `json:"typeMeta"`

	// Aggregate information about pods belonging to this Deployment
	Pods common.PodInfo `json:"pods"`

	// Container images of this Deployment
	ContainerImages []string `json:"containerImages"`

	// Init Container images of this Deployment
	InitContainerImages []string `json:"initContainerImages"`
}

func GetDeploymentList(client kubernetes.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*DeploymentList, error) {

	glog.Info("Getting list of deployments in the cluster")

	channels := &common.ResourceChannels{
		DeploymentList: common.GetDeploymentListChannel(client, nsQuery, 1),
		PodList: common.GetPodListChannel(client, nsQuery, 1),
		EventList: common.GetEventListChannel(client, nsQuery, 1),
		ReplicaSetList: common.GetReplicaSetListChannel(client, nsQuery, 1),
	}

	return GetDeploymentListFromChannels(channels, dsQuery)
}

func GetDeploymentListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery) (
	*DeploymentList, error) {
	deployments := <- channels.DeploymentList.List
	err := <- channels.DeploymentList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	pods := <- channels.PodList.List
	err = <- channels.PodList.Error
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	events := <- channels.EventList.List
	err = <- channels.EventList.Error
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	rs := <-channels.ReplicaSetList.List
	err = <-channels.ReplicaSetList.Error
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	deploymentList := toDeploymentList(deployments.Items, pods.Items, events.Items, rs.Items, nonCriticalErrors, dsQuery)
	return deploymentList, nil
}

func toDeploymentList(deployments []apps.Deployment, pods []v1.Pod, events []v1.Event, rs []apps.ReplicaSet,
	nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery) *DeploymentList {

	deploymentList := &DeploymentList{
		Deployments: make([]Deployment, 0),
		ListMeta: api.ListMeta{TotalItems: len(deployments)},
		Errors: nonCriticalErrors,
	}

	deploymentCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(deployments), dsQuery)
	deployments = fromCells(deploymentCells)
	deploymentList.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, deployment := range deployments {
		matchingPods := common.FilterDeploymentPodsByOwnerReference(deployment, rs, pods)
		podInfo := common.GetPodInfo(deployment.Status.Replicas, deployment.Spec.Replicas, matchingPods)
		podInfo.Warnings = event.GetPodsEventWarnings(events, matchingPods)

		deploymentList.Deployments = append(deploymentList.Deployments,
			Deployment{
				ObjectMeta:          api.NewObjectMeta(deployment.ObjectMeta),
				TypeMeta:            api.NewTypeMeta(api.ResourceKindDeployment),
				ContainerImages:     common.GetContainerImages(&deployment.Spec.Template.Spec),
				InitContainerImages: common.GetInitContainerImages(&deployment.Spec.Template.Spec),
				Pods:                podInfo,
			})
	}

	return deploymentList
}


