package controller

import (
	"fmt"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	apps "k8s.io/api/apps/v1beta2"
	batchV1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"strings"
)

// ResourceOwner is an structure representing resource owner, it may be Replication Controller,
// Daemon Set, Job etc.
type ResourceOwner struct {
	ObjectMeta          api.ObjectMeta `json:"objectMeta"`
	TypeMeta            api.TypeMeta   `json:"typeMeta"`
	Pods                common.PodInfo `json:"pods"`
	ContainerImages     []string       `json:"containerImages"`
	InitContainerImages []string       `json:"initContainerImages"`
}

// LogSources is a structure that represents all log files (all combinations of pods and container)
// from a higher level controller (such as ReplicaSet).
type LogSources struct {
	ContainerNames     []string `json:"containerNames"`
	InitContainerNames []string `json:"initContainerNames"`
	PodNames           []string `json:"podNames"`
}

// ResourceController is an interface, that allows to perform operations on resource controller
type ResourceController interface {
	// UID returns UID of controlled resource
	UID() types.UID
	// Get is a method, that returns ResourceOwner object
	Get(allPods []v1.Pod, allEvents []v1.Event) ResourceOwner
	// Returns all log sources of controlled resource
	GetLogSources(allPods []v1.Pod) LogSources
}

// NewResourceController creates instance of ResourceController based on given reference
// It allows to convert reference of owner/creator to real object
func NewResourceController(ref metaV1.OwnerReference, namespace string, client kubernetes.Interface) (
	ResourceController, error) {
	switch strings.ToLower(ref.Kind	) {
	case api.ResourceKindJob:
		job, err := client.BatchV1().Jobs(namespace).Get(ref.Name, metaV1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return JobController(*job), nil
	case api.ResourceKindPod:
		pod, err := client.CoreV1().Pods(namespace).Get(ref.Name, metaV1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return PodController(*pod), nil
	case api.ResourceKindReplicaSet:
		rs, err := client.AppsV1beta2().ReplicaSets(namespace).Get(ref.Name, metaV1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return ReplicaSetController(*rs), nil
	case api.ResourceKindReplicationController:
		rc, err := client.CoreV1().ReplicationControllers(namespace).Get(ref.Name, metaV1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return ReplicationControllerController(*rc), nil
	case api.ResourceKindDaemonSet:
		ds, err := client.AppsV1beta2().DaemonSets(namespace).Get(ref.Name, metaV1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return DaemonSetController(*ds), nil
	case api.ResourceKindStatefulSet:
		ss, err := client.AppsV1beta2().StatefulSets(namespace).Get(ref.Name, metaV1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return StatefulSetController(*ss), nil
	default:
		return nil, fmt.Errorf("Unknown reference kind %s", ref.Kind)
	}
}

type JobController batchV1.Job

func (self JobController) UID() types.UID {
	return batchV1.Job(self).UID
}

func (self JobController) Get(allPods []v1.Pod, allEvents []v1.Event) ResourceOwner {
	matchingPods := common.FilterPodsForJob(batchV1.Job(self), allPods)
	podInfo := common.GetPodInfo(self.Status.Active, self.Spec.Completions, matchingPods)
	podInfo.Warnings = event.GetPodsEventWarnings(allEvents, matchingPods)

	return ResourceOwner{
		TypeMeta:            api.NewTypeMeta(api.ResourceKindJob),
		ObjectMeta:          api.NewObjectMeta(self.ObjectMeta),
		Pods:                podInfo,
		ContainerImages:     common.GetContainerImages(&self.Spec.Template.Spec),
		InitContainerImages: common.GetInitContainerImages(&self.Spec.Template.Spec),
	}
}

func (self JobController) GetLogSources(allPods []v1.Pod) LogSources {
	controlledPods := common.FilterPodsForJob(batchV1.Job(self), allPods)
	return LogSources{
		PodNames:           getPodNames(controlledPods),
		ContainerNames:     common.GetContainerNames(&self.Spec.Template.Spec),
		InitContainerNames: common.GetInitContainerNames(&self.Spec.Template.Spec),
	}
}

type PodController v1.Pod

func (self PodController) UID() types.UID {
	return v1.Pod(self).UID
}

func (self PodController) Get(allPods []v1.Pod, allEvents []v1.Event) ResourceOwner {
	matchingPods := common.FilterPodsByControllerRef(&self, allPods)
	podInfo := common.GetPodInfo(int32(len(matchingPods)), nil, matchingPods) // Pods should not desire any Pods
	podInfo.Warnings = event.GetPodsEventWarnings(allEvents, matchingPods)

	return ResourceOwner{
		TypeMeta:            api.NewTypeMeta(api.ResourceKindPod),
		ObjectMeta:          api.NewObjectMeta(self.ObjectMeta),
		Pods:                podInfo,
		ContainerImages:     common.GetNonduplicateContainerImages(matchingPods),
		InitContainerImages: common.GetNonduplicateInitContainerImages(matchingPods),

	}
}

func (self PodController) GetLogSources(allPods []v1.Pod) LogSources {
	controlledPods := common.FilterPodsByControllerRef(&self, allPods)
	return LogSources{
		PodNames:           getPodNames(controlledPods),
		ContainerNames:     common.GetNonduplicateContainerNames(controlledPods),
		InitContainerNames: common.GetNonduplicateInitContainerNames(controlledPods),
	}
}

type ReplicaSetController apps.ReplicaSet

func (self ReplicaSetController) UID() types.UID {
	return apps.ReplicaSet(self).UID
}

func (self ReplicaSetController) Get(allPods []v1.Pod, allEvents []v1.Event) ResourceOwner {
	matchingPods := common.FilterPodsByControllerRef(&self, allPods)
	podInfo := common.GetPodInfo(self.Status.Replicas, self.Spec.Replicas, matchingPods)
	podInfo.Warnings = event.GetPodsEventWarnings(allEvents, matchingPods)

	return ResourceOwner{
		TypeMeta:            api.NewTypeMeta(api.ResourceKindReplicaSet),
		ObjectMeta:          api.NewObjectMeta(self.ObjectMeta),
		Pods:                podInfo,
		ContainerImages:     common.GetContainerImages(&self.Spec.Template.Spec),
		InitContainerImages: common.GetInitContainerImages(&self.Spec.Template.Spec),
	}
}

func (self ReplicaSetController) GetLogSources(allPods []v1.Pod) LogSources {
	controlledPods := common.FilterPodsByControllerRef(&self, allPods)
	return LogSources{
		PodNames:           getPodNames(controlledPods),
		ContainerNames:     common.GetContainerNames(&self.Spec.Template.Spec),
		InitContainerNames: common.GetInitContainerNames(&self.Spec.Template.Spec),
	}
}

type ReplicationControllerController v1.ReplicationController

func (self ReplicationControllerController) UID() types.UID {
	return v1.ReplicationController(self).UID
}

func (self ReplicationControllerController) Get(allPods []v1.Pod, allEvents []v1.Event) ResourceOwner {
	matchingPods := common.FilterPodsByControllerRef(&self, allPods)
	podInfo := common.GetPodInfo(self.Status.Replicas, self.Spec.Replicas, matchingPods)
	podInfo.Warnings = event.GetPodsEventWarnings(allEvents, matchingPods)

	return ResourceOwner{
		TypeMeta:            api.NewTypeMeta(api.ResourceKindReplicationController),
		ObjectMeta:          api.NewObjectMeta(self.ObjectMeta),
		Pods:                podInfo,
		ContainerImages:     common.GetContainerImages(&self.Spec.Template.Spec),
		InitContainerImages: common.GetInitContainerImages(&self.Spec.Template.Spec),
	}
}

func (self ReplicationControllerController) GetLogSources(allPods []v1.Pod) LogSources {
	controlledPods := common.FilterPodsByControllerRef(&self, allPods)
	return LogSources{
		PodNames:           getPodNames(controlledPods),
		ContainerNames:     common.GetContainerNames(&self.Spec.Template.Spec),
		InitContainerNames: common.GetInitContainerNames(&self.Spec.Template.Spec),
	}
}

type DaemonSetController apps.DaemonSet

func (self DaemonSetController) UID() types.UID {
	return apps.DaemonSet(self).UID
}

func (self DaemonSetController) Get(allPods []v1.Pod, allEvents []v1.Event) ResourceOwner {
	matchingPods := common.FilterPodsByControllerRef(&self, allPods)
	podInfo := common.GetPodInfo(self.Status.CurrentNumberScheduled,
		&self.Status.DesiredNumberScheduled, matchingPods)
	podInfo.Warnings = event.GetPodsEventWarnings(allEvents, matchingPods)

	return ResourceOwner{
		TypeMeta:            api.NewTypeMeta(api.ResourceKindDaemonSet),
		ObjectMeta:          api.NewObjectMeta(self.ObjectMeta),
		Pods:                podInfo,
		ContainerImages:     common.GetContainerImages(&self.Spec.Template.Spec),
		InitContainerImages: common.GetInitContainerImages(&self.Spec.Template.Spec),
	}
}

func (self DaemonSetController) GetLogSources(allPods []v1.Pod) LogSources {
	controlledPods := common.FilterPodsByControllerRef(&self, allPods)
	return LogSources{
		PodNames:           getPodNames(controlledPods),
		ContainerNames:     common.GetContainerNames(&self.Spec.Template.Spec),
		InitContainerNames: common.GetInitContainerNames(&self.Spec.Template.Spec),
	}
}

type StatefulSetController apps.StatefulSet

func (self StatefulSetController) UID() types.UID {
	return apps.StatefulSet(self).UID
}

func (self StatefulSetController) Get(allPods []v1.Pod, allEvents []v1.Event) ResourceOwner {
	matchingPods := common.FilterPodsByControllerRef(&self, allPods)
	podInfo := common.GetPodInfo(self.Status.Replicas, self.Spec.Replicas, matchingPods)
	podInfo.Warnings = event.GetPodsEventWarnings(allEvents, matchingPods)

	return ResourceOwner{
		TypeMeta:            api.NewTypeMeta(api.ResourceKindStatefulSet),
		ObjectMeta:          api.NewObjectMeta(self.ObjectMeta),
		Pods:                podInfo,
		ContainerImages:     common.GetContainerImages(&self.Spec.Template.Spec),
		InitContainerImages: common.GetInitContainerImages(&self.Spec.Template.Spec),
	}
}

func (self StatefulSetController) GetLogSources(allPods []v1.Pod) LogSources {
	controlledPods := common.FilterPodsByControllerRef(&self, allPods)
	return LogSources{
		PodNames:           getPodNames(controlledPods),
		ContainerNames:     common.GetContainerNames(&self.Spec.Template.Spec),
		InitContainerNames: common.GetInitContainerNames(&self.Spec.Template.Spec),
	}
}

func getPodNames(pods []v1.Pod) []string {
	names := make([]string, len(pods))
	for _, pod := range pods {
		names = append(names, pod.Name)
	}
	return names
}