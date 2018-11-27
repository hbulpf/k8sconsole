package deployment

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/pod"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetDeploymentPods returns list of pods targeting deployment
func GetDeploymentPods(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery,
	namespace, deploymentName string) (*pod.PodList, error) {

	deployment, err := client.AppsV1beta2().Deployments(namespace).Get(deploymentName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	channels := &common.ResourceChannels{
		PodList: common.GetPodListChannel(client, common.NewOneNamespaceQuery(namespace), 1),
		ReplicaSetList: common.GetReplicaSetListChannel(client, common.NewOneNamespaceQuery(namespace), 1),
	}

	rawPods := <- channels.PodList.List
	if err := <- channels.PodList.Error; err != nil {
		return pod.EmptyPodList, err
	}

	rawRs := <- channels.ReplicaSetList.List
	err = <- channels.ReplicaSetList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return pod.EmptyPodList, criticalError
	}

	pods := common.FilterDeploymentPodsByOwnerReference(*deployment, rawRs.Items, rawPods.Items)
	events, err := event.GetPodsEvents(client, namespace, pods)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return pod.EmptyPodList, criticalError
	}

	podList := pod.ToPodList(pods, events, nonCriticalErrors, dsQuery)
	return &podList, nil
}
