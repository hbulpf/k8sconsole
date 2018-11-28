package cronjob


import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CronJobList contains a list of CronJobs in the cluster.
type CronJobList struct {
	ListMeta          api.ListMeta       `json:"listMeta"`
	Items             []CronJob          `json:"items"`

	// Basic information about resources status on the list.
	Status common.ResourceStatus `json:"status"`

	// List of non-critical errors, that occurred during resource retrieval.
	Errors []error `json:"errors"`
}

// CronJob is a presentation layer view of Kubernetes Cron Job resource.
type CronJob struct {
	ObjectMeta   api.ObjectMeta `json:"objectMeta"`
	TypeMeta     api.TypeMeta   `json:"typeMeta"`
	Schedule     string         `json:"schedule"`
	Suspend      *bool          `json:"suspend"`
	Active       int            `json:"active"`
	LastSchedule *metav1.Time   `json:"lastSchedule"`
}

// GetCronJobList returns a list of all CronJobs in the cluster.
func GetCronJobList(client kubernetes.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*CronJobList, error) {
	glog.Info("Getting list of all cron jobs in the cluster")

	channels := &common.ResourceChannels{
		CronJobList: common.GetCronJobListChannel(client, nsQuery, 1),
	}

	return GetCronJobListFromChannels(channels, dsQuery)
}

// GetCronJobListFromChannels returns a list of all CronJobs in the cluster reading required resource
// list once from the channels.
func GetCronJobListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery) (
	*CronJobList, error) {

	cronJobs := <-channels.CronJobList.List
	err := <-channels.CronJobList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	cronJobList := toCronJobList(cronJobs.Items, nonCriticalErrors, dsQuery)
	cronJobList.Status = getStatus(cronJobs)
	return cronJobList, nil
}

func toCronJobList(cronJobs []v1beta1.CronJob, nonCriticalErrors []error,
	dsQuery *dataselect.DataSelectQuery) *CronJobList {

	list := &CronJobList{
		Items:    make([]CronJob, 0),
		ListMeta: api.ListMeta{TotalItems: len(cronJobs)},
		Errors:   nonCriticalErrors,
	}

	cronJobCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(cronJobs), dsQuery)
	cronJobs = fromCells(cronJobCells)
	list.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, cronJob := range cronJobs {
		list.Items = append(list.Items, toCronJob(&cronJob))
	}

	return list
}

func toCronJob(cj *v1beta1.CronJob) CronJob {
	return CronJob{
		ObjectMeta:   api.NewObjectMeta(cj.ObjectMeta),
		TypeMeta:     api.NewTypeMeta(api.ResourceKindCronJob),
		Schedule:     cj.Spec.Schedule,
		Suspend:      cj.Spec.Suspend,
		Active:       len(cj.Status.Active),
		LastSchedule: cj.Status.LastScheduleTime,
	}
}
