package storageclass

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	storage "k8s.io/api/storage/v1"
	"k8s.io/client-go/kubernetes"
)

// StorageClassList holds a list of storage class objects in the clu
type StorageClassList struct {
	ListMeta api.ListMeta `json:"listMeta"`
	StorageClasses []StorageClass `json:"storageClasses"`

	// List of non-critical errors, that occurred during resource retrieval
	Errors []error `json:"errors"`
}

// StorageClass is a representation of a kubernetes StorageClass object.
type StorageClass struct {
	ObjectMeta  api.ObjectMeta    `json:"objectMeta"`
	TypeMeta    api.TypeMeta      `json:"typeMeta"`

	// Provisioner is the driver expected to handle this StorageClass.
	// For example: "kubernetes.io/gce-pd" or "kubernetes.io/aws-ebs".
	// This value may not be empty
	Provisioner string            `json:"provisioner"`

	// Parameters holds parameters for the provisioner.
	Parameters  map[string]string `json:"parameters"`
}

func GetStorageClassList(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery) (
	*StorageClassList, error) {
	glog.Info("Getting list of storage classes in the cluster")

	channels := &common.ResourceChannels{
		StorageClassList: common.GetStorageClassListChannel(client, 1),
	}

	return GetStorageClassListFromChannels(channels, dsQuery)
}

func GetStorageClassListFromChannels(channels *common.ResourceChannels,
	dsQuery *dataselect.DataSelectQuery) (*StorageClassList, error) {

	storageClassList := <- channels.StorageClassList.List
	err := <- channels.StorageClassList.Error

	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return toStorageClassList(storageClassList.Items, nonCriticalErrors, dsQuery), nil
}

func toStorageClassList(storageClasses []storage.StorageClass,
	nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery) *StorageClassList {

	storageClassList := &StorageClassList{
		StorageClasses: make([]StorageClass, 0),
		ListMeta: api.ListMeta{TotalItems: len(storageClasses)},
		Errors: nonCriticalErrors,
	}

	storageClassCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(storageClasses), dsQuery)
	storageClasses = fromCells(storageClassCells)
	storageClassList.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, storageClass := range storageClasses {
		storageClassList.StorageClasses = append(storageClassList.StorageClasses, toStorageClass(&storageClass))
	}

	return storageClassList
}