package overview

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/config"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/discovery"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/workload"
	"k8s.io/client-go/kubernetes"
)

// Overview is a list of objects present in a given namespace
type Overview struct {
	config.Config `json:",inline"`
	discovery.Discovery `json:",inline"`
	workload.Workloads `json:",inline"`

	// List of non-critical errors, that occurred during resource retrieval
	Errors []error `json:"errors"`
}

func GetOverview(client kubernetes.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*Overview, error) {

	configResources, err := config.GetConfig(client, nsQuery, dsQuery)
	if err != nil {
		return &Overview{}, err
	}

	discoveryResources, err := discovery.GetDiscovery(client, nsQuery, dsQuery)
	if err != nil {
		return &Overview{}, err
	}

	workloadsResources, err := workload.GetWorkloads(client, nsQuery, dsQuery)
	if err != nil {
		return &Overview{}, err
	}

	return &Overview{
		Config: config.Config{
			ConfigMapList: configResources.ConfigMapList,
			PersistentVolumeClaimList: configResources.PersistentVolumeClaimList,
			SecretList: configResources.SecretList,
		},

		Discovery: discovery.Discovery{
			ServiceList: discoveryResources.ServiceList,
			IngressList: discoveryResources.IngressList,
		},

		Workloads: workload.Workloads{
			DeploymentList:            workloadsResources.DeploymentList,
			ReplicaSetList:            workloadsResources.ReplicaSetList,
			CronJobList:               workloadsResources.CronJobList,
			JobList:                   workloadsResources.JobList,
			ReplicationControllerList: workloadsResources.ReplicationControllerList,
			PodList:                   workloadsResources.PodList,
			DaemonSetList:             workloadsResources.DaemonSetList,
			StatefulSetList:           workloadsResources.StatefulSetList,
		},

		Errors: errors.MergeErrors(configResources.Errors, discoveryResources.Errors,
			workloadsResources.Errors),
	}, nil
}