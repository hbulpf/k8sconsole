package replicationcontroller

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/pod"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func GetReplicationControllerPods(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery,
	rcName, namespace string) (*pod.PodList, error) {
	glog.Info("Getting replication controller %s pods in namespace %s", rcName, namespace)

	pods, err := getRawReplicationControllerPods(client, rcName, namespace)
	if err != nil {
		return nil, err
	}

	events, err := event.GetPodsEvents(client, namespace, pods)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	podList := pod.ToPodList(pods, events, nonCriticalErrors, dsQuery)
	return &podList, nil
}

func getRawReplicationControllerPods(client kubernetes.Interface, rcName, namespace string) ([]v1.Pod, error) {
	rc, err := client.CoreV1().ReplicationControllers(namespace).Get(rcName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	channels := &common.ResourceChannels{
		PodList: common.GetPodListChannel(client, common.NewOneNamespaceQuery(namespace), 1),
	}

	podList := <- channels.PodList.List
	if err := <- channels.PodList.Error; err != nil {
		return nil, err
	}

	return common.FilterPodsByControllerRef(rc, podList.Items), nil
}

// getReplicationControllerPodInfo returns simple info about pods(running, desired, failing, etc)
func getReplicationControllerPodInfo(client kubernetes.Interface, rc *v1.ReplicationController, namespace string) (
	*common.PodInfo, error) {

	labelSelector := labels.SelectorFromSet(rc.Spec.Selector)
	channels := &common.ResourceChannels{
		PodList: common.GetPodListChannelWithOptions(client, common.NewOneNamespaceQuery(namespace),
			metaV1.ListOptions{
				LabelSelector: labelSelector.String(),
				FieldSelector: fields.Everything().String(),
			}, 1),
	}

	pods := <- channels.PodList.List
	if err := <- channels.PodList.Error; err != nil {
		return nil, err
	}

	podInfo := common.GetPodInfo(rc.Status.Replicas, rc.Spec.Replicas, pods.Items)
	return &podInfo, nil
}