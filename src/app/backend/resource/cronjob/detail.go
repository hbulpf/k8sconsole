package cronjob

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/job"
	batch2 "k8s.io/api/batch/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CronJobDetail contains Cron Job details.
type CronJobDetail struct {
	ConcurrencyPolicy       string           `json:"concurrencyPolicy"`
	StartingDeadLineSeconds *int64           `json:"startingDeadlineSeconds"`
	ActiveJobs              job.JobList      `json:"activeJobs"`
	InactiveJobs            job.JobList      `json:"inactiveJobs"`
	Events                  common.EventList `json:"events"`

	// Extends list item structure.
	CronJob `json:",inline"`

	// List of non-critical errors, that occurred during resource retrieval.
	Errors []error `json:"errors"`
}

// GetCronJobDetail gets Cron Job details.
func GetCronJobDetail(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery, namespace, name string) (*CronJobDetail, error) {

	rawObject, err := client.BatchV1beta1().CronJobs(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	activeJobs, err := GetCronJobJobs(client, dsQuery, namespace, name)
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	inactiveJobs, err := GetCronJobCompletedJobs(client, dsQuery, namespace, name)

	events, err := GetCronJobEvents(client, dsQuery, namespace, name)
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	cj := toCronJobDetail(rawObject, *activeJobs, *inactiveJobs, *events, nonCriticalErrors)
	return &cj, nil
}

func toCronJobDetail(cj *batch2.CronJob, activeJobs job.JobList, inactiveJobs job.JobList, events common.EventList,
	nonCriticalErrors []error) CronJobDetail {
	return CronJobDetail{
		CronJob:                 toCronJob(cj),
		ConcurrencyPolicy:       string(cj.Spec.ConcurrencyPolicy),
		StartingDeadLineSeconds: cj.Spec.StartingDeadlineSeconds,
		ActiveJobs:              activeJobs,
		InactiveJobs:            inactiveJobs,
		Events:                  events,
		Errors:                  nonCriticalErrors,
	}
}

