package common

import (
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Condition represents a single condition of a node or pod.
// e.g. v1.Pod.Status.Condition
type Condition struct {
	// Type of condition
	Type string `json:"type"`

	// Status of condition
	Status v1.ConditionStatus `json:"status"`

	// Last probe time of a condition
	LastProbeTime metaV1.Time `json:"lastProbeTime"`

	// Last transition time of a condition
	LastTransitionTime metaV1.Time `json:"lastTransitionTime"`

	// Reason of a condition
	Reason string `json:"reason"`

	// Message of a condition.
	Message string `json:"message"`
}
