package common

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EventList struct {
	ListMeta api.ListMeta	`json:"listMeta"`
	Events []Event `json:"events"`
}

type Event struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta api.TypeMeta `json:"typeMeta"`

	// a human-readable description of the status of related object.
	Message string `json:"message"`

	// Component from which the event is generated.
	SourceComponent string `json:"sourceComponent"`

	// Host name on which the event is generated.
	SourceHost string	`json:"sourceHost"`

	// An object triggered an event.
	SubObject string `json:"object"`

	// The number of times this event has occurred.
	Count int32 `json:"count"`

	// The time at which the event was first occurred.
	FirstSeen v1.Time `json:"firstSeen"`

	// The time at which the event was last occurred.
	LastSeen v1.Time `json:"lastSeen"`

	// Short, machine-understandable string that gives the
	// reason for this event being generated
	Reason string `json:"reason"`

	// Event type
	Type string `json:"type"`
}
