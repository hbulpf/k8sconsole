package common

import "k8s.io/api/core/v1"

type PodInfo struct {
	// Number of pods that are created
	Current int32 `json:"current"`

	// Number of pods that are desired
	Desired *int32 `json:"desired"`

	// Number of pods that are currently running
	Running int32 `json:"running"`

	// Number of pods that are currently pending
	Pending int32 `json:"pending"`

	// Number of pods that are failed
	Failed int32 `json:"failed"`

	// Number of pods that are succeed
	Succeeded int32 `json:"succeeded"`

	// Unique warning messages related to pods in this resource.
	Warnings []Event `json:"warnings"`
}

// GetPodInfo returns aggregate information about a group of pods.
func GetPodInfo(current int32, desired *int32, pods []v1.Pod) PodInfo {
	result := PodInfo{
		Current: current,
		Desired: desired,
		Warnings: make([]Event, 0),
	}

	for _, pod := range pods {
		switch pod.Status.Phase {
		case v1.PodRunning:
			result.Running++
		case v1.PodPending:
			result.Pending++
		case v1.PodFailed:
			result.Failed++
		case v1.PodSucceeded:
			result.Succeeded++
		}
	}

	return result
}
