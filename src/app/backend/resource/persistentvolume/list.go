package persistentvolume

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// PersistentVolumeList contains a list of Persistent Volumes in the cluster.
type PersistentVolumeList struct {
	ListMeta api.ListMeta       `json:"listMeta"`
	Items    []PersistentVolume `json:"items"`

	// list of non-critical errors, that occurred during resource retrieval.
	Errors   []error            `json:"errors"`
}

// PersistentVolume provides the simplified presentation layer view of kubernetes Persistent Volume resource.
type PersistentVolume struct {
	ObjectMeta    api.ObjectMeta                   `json:"objectMeta"`
	TypeMeta      api.TypeMeta                     `json:"typeMeta"`
	Capacity      v1.ResourceList                  `json:"capacity"`
	AccessModes   []v1.PersistentVolumeAccessMode  `json:"accessModes"`
	ReclaimPolicy v1.PersistentVolumeReclaimPolicy `json:"reclaimPolicy"`
	StorageClass  string                           `json:"storageClass"`
	Status        v1.PersistentVolumePhase         `json:"status"`
	Claim         string                           `json:"claim"`
	Reason        string                           `json:"reason"`
}

// GetPersistentVolumeList returns a list of all Persistent Volumes in the cluster.
func GetPersistentVolumeList(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery) (
	*PersistentVolumeList, error) {
	glog.Info("Getting list persistent volumes")
	channels := &common.ResourceChannels{
		PersistentVolumeList: common.GetPersistentVolumeListChannel(client, 1),
	}
	return GetPersistentVolumeListFromChannels(channels, dsQuery)
}

// GetPersistentVolumeListFromChannels returns a list of all Persistent Volumes in the cluster
// reading required resource list once from the channels.
func GetPersistentVolumeListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery) (
	*PersistentVolumeList, error) {
	persistentVolumes := <- channels.PersistentVolumeList.List
	err := <- channels.PersistentVolumeList.Error

	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}
	return toPersistentVolumeList(persistentVolumes.Items, nonCriticalErrors, dsQuery), nil
}

func toPersistentVolumeList(persistentVolumes []v1.PersistentVolume, nonCriticalErrors []error,
	dsQuery *dataselect.DataSelectQuery) *PersistentVolumeList {

	result := &PersistentVolumeList{
		ListMeta: api.ListMeta{TotalItems: len(persistentVolumes)},
		Items: make([]PersistentVolume, 0),
		Errors: nonCriticalErrors,
	}

	pvCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(persistentVolumes), dsQuery)
	persistentVolumes = fromCells(pvCells)
	result.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, item := range persistentVolumes {
		result.Items = append(result.Items, toPersistentVolume(item))
	}

	return result
}

func toPersistentVolume(pv v1.PersistentVolume) PersistentVolume {
	return PersistentVolume{
		ObjectMeta: api.NewObjectMeta(pv.ObjectMeta),
		TypeMeta: api.NewTypeMeta(api.ResourceKindPersistentVolume),
		Capacity: pv.Spec.Capacity,
		AccessModes: pv.Spec.AccessModes,
		ReclaimPolicy: pv.Spec.PersistentVolumeReclaimPolicy,
		StorageClass: pv.Spec.StorageClassName,
		Status: pv.Status.Phase,
		Claim: getPersistentVolumeClaim(&pv),
		Reason: pv.Status.Reason,
	}
}