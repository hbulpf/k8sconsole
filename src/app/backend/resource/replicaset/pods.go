package replicaset

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/pod"
	apps "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func GetReplicaSetPods(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery,
	petSetName, namespace string) (*pod.PodList, error) {

	pods, err := getRawReplicaSetPods(client, petSetName, namespace)
	if err != nil {
		return pod.EmptyPodList, err
	}

	events, err := event.GetPodsEvents(client, namespace, pods)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	podList := pod.ToPodList(pods, events, nonCriticalErrors, dsQuery)
	return &podList, nil
}

func getRawReplicaSetPods(client kubernetes.Interface, petSetName, namespace string) ([]v1.Pod, error) {
	rs, err := client.AppsV1beta2().ReplicaSets(namespace).Get(petSetName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	channels := &common.ResourceChannels{
		PodList: common.GetPodListChannel(client, common.NewOneNamespaceQuery(namespace), 1),
	}

	podList := <- channels.PodList.List
	if err := <-channels.PodList.Error; err != nil {
		return nil, err
	}

	return common.FilterPodsByControllerRef(rs, podList.Items), nil
}

func getReplicaSetPodInfo(client kubernetes.Interface, replicaSet *apps.ReplicaSet) (*common.PodInfo, error) {
	labelselector := labels.SelectorFromSet(replicaSet.Spec.Selector.MatchLabels)
	channels := &common.ResourceChannels{
		PodList: common.GetPodListChannelWithOptions(client, common.NewOneNamespaceQuery(replicaSet.Namespace),
			metaV1.ListOptions{
				LabelSelector: labelselector.String(),
				FieldSelector: fields.Everything().String(),
			}, 1),
	}

	pods := <- channels.PodList.List
	if err := <- channels.PodList.Error; err != nil {
		return nil, err
	}

	podInfo := common.GetPodInfo(replicaSet.Status.Replicas, replicaSet.Spec.Replicas, pods.Items)
	return &podInfo, nil
}
