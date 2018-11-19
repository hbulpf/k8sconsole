package event

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"strings"
)

// FailedReasonPartials  is an array of partial strings to correctly filter warning events.
var FailedReasonPartials = []string{"failed", "err", "exceeded", "invalid", "unhealthy",
	"mismatch", "insufficient", "conflict", "outof", "nil", "backoff"}

// GetPodsEventWarnings returns warning pod events by filtering out events targeting only given pods
func GetPodsEventWarnings(events []v1.Event, pods []v1.Pod) []common.Event {
	result := make([]common.Event, 0)

	// filter out only warning events
	events = getWarningEvents(events)

	// filter out ready and successful pods
	failedPods := make([]v1.Pod, 0)
	for _, pod := range pods {
		if !isReadyOrSucceeded(pod) {
			failedPods = append(failedPods, pod)
		}
	}

	events = filterEventsByPodsUID(events, failedPods)
	events = removeDuplicates(events)

	for _, event := range events {
		result = append(result, common.Event{
			Message: event.Message,
			Reason: event.Reason,
			Type: event.Type,
		})
	}

	return result
}

// Returns filtered list of event objects. Events list is filtered to get only events targeting
// pods on the list.
func filterEventsByPodsUID(events []v1.Event, pods []v1.Pod) []v1.Event {
	result := make([]v1.Event, 0)
	podEventMap := make(map[types.UID]bool, 0)

	if len(pods) == 0 || len(events) == 0 {
		return result
	}

	for _, pod := range pods {
		podEventMap[pod.UID] = true
	}

	for _, event := range events {
		if _, exists := podEventMap[event.InvolvedObject.UID]; exists {
			result = append(result, event)
		}
	}

	return result
}

// Returns filtered list of event objects.
// Event list object is filtered to get only warning events.
func getWarningEvents(events []v1.Event) []v1.Event {
	return filterEventsByType(events, v1.EventTypeWarning)
}

// Filters kubernetes API event objects based on event type.
// Empty string will return all events.
func filterEventsByType(events []v1.Event, eventType string) []v1.Event {
	if len(events) == 0 || len(eventType) == 0 {
		return events
	}

	result := make([]v1.Event, 0)
	for _, event := range events {
		if event.Type == eventType {
			result = append(result, event)
		}
	}

	return result
}

// Returns true if reason string contains any partial string indicating that this may be a
// warning, false otherwise
func isFailedReason(reason string, partials ...string) bool {
	for _, partial := range partials {
		if strings.Contains(strings.ToLower(reason), partial) {
			return true
		}
	}

	return false
}

func removeDuplicates(slice []v1.Event) []v1.Event {
	visited := make(map[string]bool, 0)
	result := make([]v1.Event, 0)

	for _, elem := range slice {
		if !visited[elem.Reason] {
			visited[elem.Reason] = true
			result = append(result, elem)
		}
	}

	return result
}

// Returns true if given pod is in state ready or succeeded, false otherwise
func isReadyOrSucceeded(pod v1.Pod) bool {
	if pod.Status.Phase == v1.PodSucceeded {
		return true
	}

	if pod.Status.Phase == v1.PodRunning {
		for _, c := range pod.Status.Conditions {
			if c.Type == v1.PodReady &&
				 c.Status == v1.ConditionFalse {
				return false
			}
		}

		return true
	}

	return false
}