package deployment

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/replicaset"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apps "k8s.io/api/apps/v1beta2"
)

func GetDeploymentOldReplicaSets(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery,
	namespace, deploymentName string) (*replicaset.ReplicaSetList, error) {

	oldReplicaSetList := &replicaset.ReplicaSetList{
		ReplicaSets: make([]replicaset.ReplicaSet, 0),
		ListMeta: api.ListMeta{TotalItems: 0},
	}

	deployment, err := client.AppsV1beta2().Deployments(namespace).Get(deploymentName, metaV1.GetOptions{})
	if err != nil {
		return oldReplicaSetList, err
	}

	selector, err := metaV1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return oldReplicaSetList, err
	}
	options := metaV1.ListOptions{LabelSelector: selector.String()}

	channels := &common.ResourceChannels{
		ReplicaSetList: common.GetReplicaSetListChannelWithOptions(client,
			common.NewOneNamespaceQuery(namespace), options, 1),
		PodList: common.GetPodListChannelWithOptions(client,
			common.NewOneNamespaceQuery(namespace), options, 1),
		EventList: common.GetEventListChannelWithOptions(client,
			common.NewOneNamespaceQuery(namespace), options, 1),
	}

	rawRs := <- channels.ReplicaSetList.List
	if err := <- channels.ReplicaSetList.Error; err != nil {
		return oldReplicaSetList, err
	}

	rawPods := <- channels.PodList.List
	if err := <- channels.PodList.Error; err != nil {
		return oldReplicaSetList, err
	}

	rawEvents := <-channels.EventList.List
	err = <-channels.EventList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return oldReplicaSetList, criticalError
	}

	rawRepSets := make([]*apps.ReplicaSet, 0)
	for i := range rawRs.Items {
		rawRepSets = append(rawRepSets, &rawRs.Items[i])
	}
	oldRs, _, err := FindOldReplicaSets(deployment, rawRepSets)
	if err != nil {
		return oldReplicaSetList, err
	}

	oldReplicaSets := make([]apps.ReplicaSet, len(oldRs))
	for i, replicaSet := range oldRs {
		oldReplicaSets[i] = *replicaSet
	}

	oldReplicaSetList = replicaset.ToReplicaSetList(oldReplicaSets,
		rawPods.Items, rawEvents.Items, nonCriticalErrors, dsQuery)

	return oldReplicaSetList, nil
}
