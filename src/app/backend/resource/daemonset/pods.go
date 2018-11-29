package daemonset

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/pod"
	apps "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetDaemonSetPods return list of pods targeting daemon set.
func GetDaemonSetPods(client kubernetes.Interface,
	dsQuery *dataselect.DataSelectQuery, daemonSetName, namespace string) (*pod.PodList, error) {
	glog.Infof("Getting replication controller %s pods in namespace %s", daemonSetName, namespace)

	pods, err := getRawDaemonSetPods(client, daemonSetName, namespace)
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

// Returns array of api pods targeting daemon set with given name.
func getRawDaemonSetPods(client kubernetes.Interface, daemonSetName, namespace string) ([]v1.Pod, error) {
	daemonSet, err := client.AppsV1beta2().DaemonSets(namespace).Get(daemonSetName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	channels := &common.ResourceChannels{
		PodList: common.GetPodListChannel(client, common.NewOneNamespaceQuery(namespace), 1),
	}

	podList := <-channels.PodList.List
	if err := <-channels.PodList.Error; err != nil {
		return nil, err
	}

	matchingPods := common.FilterPodsByControllerRef(daemonSet, podList.Items)
	return matchingPods, nil
}

// Returns simple info about pods(running, desired, failing, etc.) related to given daemon set.
func getDaemonSetPodInfo(client kubernetes.Interface, daemonSet *apps.DaemonSet) (
	*common.PodInfo, error) {

	pods, err := getRawDaemonSetPods(client, daemonSet.Name, daemonSet.Namespace)
	if err != nil {
		return nil, err
	}

	podInfo := common.GetPodInfo(daemonSet.Status.CurrentNumberScheduled,
		&daemonSet.Status.DesiredNumberScheduled, pods)
	return &podInfo, nil
}
