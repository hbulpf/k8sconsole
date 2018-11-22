package persistentvolumeclaim

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	kcErrors "github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// PersistentVolumeClaimList contains a list of Persistent Volume Claims in the cluster.
type PersistentVolumeClaimList struct {
	ListMeta api.ListMeta							`json:"listMeta"`

	// Unordered list of persistent volume claim
	Items    []PersistentVolumeClaim	`json:"items"`

	// List of non-critical errors, that occurred during resource retrieval.
	Errors   []error									`json:"errors"`
}

// PersistentVolumeClaim provides the simplified presentation layer view of Kubernetes Persistent Volume Claim resource.
type PersistentVolumeClaim struct {
	ObjectMeta   api.ObjectMeta                  `json:"objectMeta"`
	TypeMeta     api.TypeMeta                    `json:"typeMeta"`
	Status       string                          `json:"status"`
	Volume       string                          `json:"volume"`
	Capacity     v1.ResourceList                 `json:"capacity"`
	AccessModes  []v1.PersistentVolumeAccessMode `json:"accessModes"`
	StorageClass *string                         `json:"storageClass"`
}

func GetPersistentVolumeClaimList(client kubernetes.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*PersistentVolumeClaimList, error) {

	glog.Info("Getting list persistent volumes claims")
	channels := &common.ResourceChannels{
		PersistentVolumeClaimList: common.GetPersistentVolumeClaimListChannel(client, nsQuery, 1),
	}

	return GetPersistentVolumeClaimListFromChannels(channels, nsQuery, dsQuery)
}

func GetPersistentVolumeClaimListFromChannels(channels *common.ResourceChannels, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*PersistentVolumeClaimList, error) {

	persistentVolumeClaims := <- channels.PersistentVolumeClaimList.List
	err := <- channels.PersistentVolumeClaimList.Error
	nonCriticalErrors, criticalError := kcErrors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return toPersistentVolumeChaimList(persistentVolumeClaims.Items, nonCriticalErrors, dsQuery), nil
}

func toPersistentVolumeChaimList(persistentVolumeClaims []v1.PersistentVolumeClaim,
	nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery) *PersistentVolumeClaimList {

	result := &PersistentVolumeClaimList{
		ListMeta: api.ListMeta{TotalItems: len(persistentVolumeClaims)},
		Items: make([]PersistentVolumeClaim, 0),
		Errors: nonCriticalErrors,
	}

	pvcCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(persistentVolumeClaims), dsQuery)
	persistentVolumeClaims = fromCells(pvcCells)
	result.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, item := range persistentVolumeClaims {
		result.Items = append(result.Items, toPersistentVolumeClaim(item))
	}

	return result
}

func toPersistentVolumeClaim(pvc v1.PersistentVolumeClaim) PersistentVolumeClaim {
	return PersistentVolumeClaim{
		ObjectMeta:   api.NewObjectMeta(pvc.ObjectMeta),
		TypeMeta:     api.NewTypeMeta(api.ResourceKindPersistentVolumeClaim),
		Status:       string(pvc.Status.Phase),
		Volume:       pvc.Spec.VolumeName,
		Capacity:     pvc.Status.Capacity,
		AccessModes:  pvc.Spec.AccessModes,
		StorageClass: pvc.Spec.StorageClassName,
	}
}
