package common

import (
	appV1 "k8s.io/api/apps/v1beta2"
	batchV1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FilterPodsByControllerResource returns a subset of pods controlled by given deployment.
func FilterDeploymentPodsByOwnerReference(deployment appV1.Deployment, allRS []appV1.ReplicaSet,
	allPods []v1.Pod) []v1.Pod {
	var matchingPods []v1.Pod
	for _, rs := range allRS {
		if metaV1.IsControlledBy(&rs, &deployment) {
			matchingPods = append(matchingPods, FilterPodsByControllerRef(&rs, allPods)...)
		}
	}

	return matchingPods
}

// FilterPodsByControllerRef returns a subset of pods controlled by given controller resource, excluding deployments.
func FilterPodsByControllerRef(owner metaV1.Object, allPods []v1.Pod) []v1.Pod {
	var matchingPods []v1.Pod
	for _, pod := range allPods {
		if metaV1.IsControlledBy(&pod, owner) {
			matchingPods = append(matchingPods, pod)
		}
	}
	return matchingPods
}

// FilterPodsForJob returns a subset of pods associated with given job
func FilterPodsForJob(job batchV1.Job, pods []v1.Pod) []v1.Pod {
	var matchingPods []v1.Pod
	for _, pod := range pods {
		if pod.Namespace == job.Namespace &&
			pod.Labels["controller-uid"] == job.Spec.Selector.MatchLabels["controller-uid"] {
			matchingPods = append(matchingPods, pod)
		}
	}

	return matchingPods
}

// GetContainerImages returns container image strings from the given pod spec.
func GetContainerImages(podTemplate *v1.PodSpec) []string {
	var containerImages []string
	for _, container := range podTemplate.Containers {
		containerImages = append(containerImages, container.Image)
	}
	return containerImages
}

// GetInitContainerImages returns init container image strings from the given pod spec.
func GetInitContainerImages(podTemplate *v1.PodSpec) []string {
	var initContainerImages []string
	for _, initContainer := range podTemplate.InitContainers {
		initContainerImages = append(initContainerImages, initContainer.Image)
	}
	return initContainerImages
}

// GetContainerNames returns the container image name without the version number from the given pod spec.
func GetContainerNames(podTemplate *v1.PodSpec) []string {
	var containerNames []string
	for _, container := range podTemplate.Containers {
		containerNames = append(containerNames, container.Name)
	}
	return containerNames
}

// GetInitContainerNames returns the init container image name without the version number from the given pod spec.
func GetInitContainerNames(podTemplate *v1.PodSpec) []string {
	var initContainerNames []string
	for _, initContainer := range podTemplate.InitContainers {
		initContainerNames = append(initContainerNames, initContainer.Name)
	}
	return initContainerNames
}

func GetNonduplicateContainerImages(pods []v1.Pod) []string {
	var containerImages []string
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			if noStringInSlice(container.Image, containerImages) {
				containerImages = append(containerImages, container.Image)
			}
		}
	}
	return containerImages
}

// GetNonduplicateInitContainerImages returns list of init container image strings without duplicates
func GetNonduplicateInitContainerImages(podList []v1.Pod) []string{
	var initContainerImages []string
	for _, pod := range podList {
		for _, initContainer := range pod.Spec.InitContainers {
			if noStringInSlice(initContainer.Image, initContainerImages){
				initContainerImages = append(initContainerImages, initContainer.Image)
			}
		}
	}
	return initContainerImages
}

// GetNonduplicateContainerNames returns list of container names strings without duplicates
func GetNonduplicateContainerNames(podList []v1.Pod) []string{
	var containerNames []string
	for _, pod := range podList {
		for _, container := range pod.Spec.Containers {
			if noStringInSlice(container.Name, containerNames){
				containerNames = append(containerNames, container.Name)
			}
		}
	}
	return containerNames
}

// GetNonduplicateInitContainerNames returns list of init container names strings without duplicates
func GetNonduplicateInitContainerNames(podList []v1.Pod) []string{
	var initContainerNames []string
	for _, pod := range podList {
		for _, initContainer := range pod.Spec.InitContainers {
			if noStringInSlice(initContainer.Name, initContainerNames){
				initContainerNames = append(initContainerNames, initContainer.Name)
			}
		}
	}
	return initContainerNames
}

//noStringInSlice checks if string in array
func noStringInSlice(str string, arr []string) bool {
	for _, alreadyStr := range arr {
		if str == alreadyStr {
			return false
		}
	}
	return true
}

// EqualIgnoreHash returns true if two given podTemplateSpec are equal, ignoring the diff in value of Labels[pod-template-hash]
func EqualIgnoreHash(template1, template2 v1.PodTemplateSpec) bool {
	labels1, labels2 := template1.Labels, template2.Labels
	if len(labels1) > len(labels2) {
		labels1, labels2 = labels2, labels1
	}

	for k, v := range labels2 {
		if labels1[k] != v && k != appV1.DefaultDaemonSetUniqueLabelKey {
			return false
		}
	}
	template1.Labels, template2.Labels = nil, nil
	return equality.Semantic.DeepEqual(template1, template2)
}
