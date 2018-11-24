package service

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/pod"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func GetServicePods(client kubernetes.Interface, namespace, name string,
	dsQuery *dataselect.DataSelectQuery) (*pod.PodList, error) {
	podList := pod.PodList{
		Pods: []pod.Pod{},
	}

	service, err := client.CoreV1().Services(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		return &podList, err
	}

	if service.Spec.Selector == nil {
		return &podList, nil
	}

	labelSelector := labels.SelectorFromSet(service.Spec.Selector)
	channels := &common.ResourceChannels{
		PodList: common.GetPodListChannelWithOptions(client, common.NewOneNamespaceQuery(namespace),
			metaV1.ListOptions{
				LabelSelector: labelSelector.String(),
				FieldSelector: fields.Everything().String(),
			}, 1),
	}

	apiPodList := <- channels.PodList.List
	if err := <- channels.PodList.Error; err != nil {
		return &podList, err
	}

	events, err := event.GetPodsEvents(client, namespace, apiPodList.Items)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return &podList, criticalError
	}

	podList = pod.ToPodList(apiPodList.Items, events, nonCriticalErrors, dsQuery)
	return &podList, nil
}