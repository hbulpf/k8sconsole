package workload

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/cronjob"
	ds "github.com/wzt3309/k8sconsole/src/app/backend/resource/daemonset"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/deployment"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/job"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/pod"
	rs "github.com/wzt3309/k8sconsole/src/app/backend/resource/replicaset"
	rc "github.com/wzt3309/k8sconsole/src/app/backend/resource/replicationcontroller"
	sts "github.com/wzt3309/k8sconsole/src/app/backend/resource/statefulset"
	"k8s.io/client-go/kubernetes"
)

// Workloads structure contains all resource lists grouped into the workloads category
type Workloads struct {
	DeploymentList            deployment.DeploymentList `json:"deploymentList"`
	ReplicaSetList            rs.ReplicaSetList `json:"replicaSetList"`
	JobList                   job.JobList `json:"jobList"`
	CronJobList               cronjob.CronJobList `json:"cronJobList"`
	ReplicationControllerList rc.ReplicationControllerList `json:"replicationControllerList"`
	PodList                   pod.PodList `json:"podList"`
	DaemonSetList             ds.DaemonSetList `json:"daemonSetList"`
	StatefulSetList           sts.StatefulSetList `json:"statefulSetList"`

	// List of non-critical errors, that occurred during resource retrieval.
	Errors []error `json:"errors"`
}

func GetWorkloads(client kubernetes.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*Workloads, error) {

	glog.Info("Getting list of all workloads")
	channels := &common.ResourceChannels{
		ReplicationControllerList: common.GetReplicationControllerListChannel(client, nsQuery, 1),
		ReplicaSetList:            common.GetReplicaSetListChannel(client, nsQuery, 2),
		JobList:                   common.GetJobListChannel(client, nsQuery, 1),
		CronJobList:               common.GetCronJobListChannel(client, nsQuery, 1),
		DeploymentList:            common.GetDeploymentListChannel(client, nsQuery, 1),
		DaemonSetList:             common.GetDaemonSetListChannel(client, nsQuery, 1),
		StatefulSetList:           common.GetStatefulSetListChannel(client, nsQuery, 1),
		ServiceList:               common.GetServiceListChannel(client, nsQuery, 1),
		PodList:                   common.GetPodListChannel(client, nsQuery, 7),
		EventList:                 common.GetEventListChannel(client, nsQuery, 7),
	}

	return GetWorkloadsFromChannels(channels, dsQuery)
}

func GetWorkloadsFromChannels(channels *common.ResourceChannels,
	dsQuery *dataselect.DataSelectQuery) (*Workloads, error) {

	numErrs := 8
	errChan := make(chan error, numErrs)
	deployChan := make(chan *deployment.DeploymentList)
	rsChan := make(chan *rs.ReplicaSetList)
	jobChan := make(chan *job.JobList)
	cronJobChan := make(chan *cronjob.CronJobList)
	rcChan := make(chan *rc.ReplicationControllerList)
	podChan := make(chan *pod.PodList)
	dsChan := make(chan *ds.DaemonSetList)
	stsChan := make(chan *sts.StatefulSetList)

	go func() {
		items, err := deployment.GetDeploymentListFromChannels(channels, dsQuery)
		errChan <- err
		deployChan <- items
	}()

	go func() {
		items, err := rs.GetReplicaSetListFromChannels(channels, dsQuery)
		errChan <- err
		rsChan <- items
	}()

	go func() {
		items, err := job.GetJobListFromChannels(channels, dsQuery)
		errChan <- err
		jobChan <- items
	}()

	go func() {
		items, err := cronjob.GetCronJobListFromChannels(channels, dsQuery)
		errChan <- err
		cronJobChan <- items
	}()

	go func() {
		items, err := rc.GetReplicationControllerListFromChannels(channels, dsQuery)
		errChan <- err
		rcChan <- items
	}()

	go func() {
		items, err := pod.GetPodListFromChannels(channels, dsQuery)
		errChan <- err
		podChan <- items
	}()

	go func() {
		items, err := ds.GetDaemonSetListFromChannels(channels, dsQuery)
		errChan <- err
		dsChan <- items
	}()

	go func() {
		items, err := sts.GetStatefulSetListFromChannels(channels, dsQuery	)
		errChan <- err
		stsChan <- items
	}()

	for i := 0; i < numErrs; i++ {
		err := <- errChan
		if err != nil {
			return nil, err
		}
	}

	workloads := &Workloads{
		DeploymentList:            *(<-deployChan),
		ReplicaSetList:            *(<-rsChan),
		JobList:                   *(<-jobChan),
		CronJobList:               *(<-cronJobChan),
		ReplicationControllerList: *(<-rcChan),
		PodList:                   *(<-podChan),
		DaemonSetList:             *(<-dsChan),
		StatefulSetList:           *(<-stsChan),
	}

	workloads.Errors = errors.MergeErrors(
		workloads.DeploymentList.Errors,
		workloads.ReplicaSetList.Errors,
		workloads.JobList.Errors,
		workloads.CronJobList.Errors,
		workloads.ReplicationControllerList.Errors,
		workloads.PodList.Errors,
		workloads.DaemonSetList.Errors,
		workloads.StatefulSetList.Errors)

	return workloads, nil
}