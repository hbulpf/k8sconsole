package job

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/pod"
	batch "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// GetJobPods return list of pods targeting job.
func GetJobPods(client kubernetes.Interface,	dsQuery *dataselect.DataSelectQuery,
	namespace string, jobName string) (*pod.PodList, error) {
	glog.Infof("Getting replication controller %s pods in namespace %s", jobName, namespace)

	pods, err := getRawJobPods(client, jobName, namespace)
	if err != nil {
		return pod.EmptyPodList, err
	}

	events, err := event.GetPodsEvents(client, namespace, pods)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return pod.EmptyPodList, criticalError
	}

	podList := pod.ToPodList(pods, events, nonCriticalErrors, dsQuery)
	return &podList, nil
}

// Returns array of api pods targeting job with given name.
func getRawJobPods(client kubernetes.Interface, petSetName, namespace string) ([]v1.Pod, error) {
	job, err := client.Batch().Jobs(namespace).Get(petSetName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	labelSelector := labels.SelectorFromSet(job.Spec.Selector.MatchLabels)
	channels := &common.ResourceChannels{
		PodList: common.GetPodListChannelWithOptions(client, common.NewOneNamespaceQuery(namespace),
			metaV1.ListOptions{
				LabelSelector: labelSelector.String(),
				FieldSelector: fields.Everything().String(),
			}, 1),
	}

	podList := <-channels.PodList.List
	if err := <-channels.PodList.Error; err != nil {
		return nil, err
	}

	return podList.Items, nil
}

// Returns simple info about pods(running, desired, failing, etc.) related to given job.
func getJobPodInfo(client kubernetes.Interface, job *batch.Job) (*common.PodInfo, error) {
	labelSelector := labels.SelectorFromSet(job.Spec.Selector.MatchLabels)
	channels := &common.ResourceChannels{
		PodList: common.GetPodListChannelWithOptions(client, common.NewOneNamespaceQuery(
			job.Namespace),
			metaV1.ListOptions{
				LabelSelector: labelSelector.String(),
				FieldSelector: fields.Everything().String(),
			}, 1),
	}

	pods := <-channels.PodList.List
	if err := <-channels.PodList.Error; err != nil {
		return nil, err
	}

	podInfo := common.GetPodInfo(job.Status.Active, job.Spec.Completions, pods.Items)

	// This pod info for jobs should be get from job status, similar to kubectl describe logic.
	podInfo.Running = job.Status.Active
	podInfo.Succeeded = job.Status.Succeeded
	podInfo.Failed = job.Status.Failed
	return &podInfo, nil
}
