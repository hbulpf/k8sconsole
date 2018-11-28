package cronjob

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	"k8s.io/client-go/kubernetes"
)

// GetCronJobEvents gets events associated to cron job.
func GetCronJobEvents(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery, namespace, name string) (
	*common.EventList, error) {

	raw, err := event.GetEvents(client, namespace, name)
	if err != nil {
		return event.EmptyEventList, err
	}

	events := event.CreateEventList(raw, dsQuery)
	return &events, nil
}
