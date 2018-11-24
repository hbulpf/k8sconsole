package pod

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	"k8s.io/client-go/kubernetes"
)

func GetEventsForPod(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery,
	namespace, name string) (*common.EventList, error) {
	eventList := common.EventList{
		Events: make([]common.Event, 0),
		ListMeta: api.ListMeta{TotalItems: 0},
	}

	podEvents, err := event.GetPodEvents(client, namespace, name)
	if err != nil {
		return &eventList, err
	}

	eventList = event.CreateEventList(podEvents, dsQuery)

	glog.Infof("Found %d events related to %s pod in %s namespace", len(eventList.Events), name, namespace)

	return &eventList, nil
}