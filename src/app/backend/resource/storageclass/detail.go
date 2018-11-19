package storageclass

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/persistentvolume"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// StorageClassDetail provides the presentation layer view of Kubernetes StorageClass resource,
// It is StorageClassDetail plus PersistentVolumes associated with StorageClass.
type StorageClassDetail struct {
	ObjectMeta           api.ObjectMeta                        `json:"objectMeta"`
	TypeMeta             api.TypeMeta                          `json:"typeMeta"`
	Provisioner          string                                `json:"provisioner"`
	Parameters           map[string]string                     `json:"parameters"`
	PersistentVolumeList persistentvolume.PersistentVolumeList `json:"persistentVolumeList"`
}

func GetStorageClassDetail(client kubernetes.Interface, name string) (*StorageClassDetail, error) {
	glog.Infof("Getting details of %s storage class", name)

	storage, err := client.StorageV1().StorageClasses().Get(name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	persistentVolumeList, err := persistentvolume.GetStorageClassPersistentVolumes(client,
		storage.Name, dataselect.DefaultDataSelect)
	storageClass := toStorageClassDetail(storage, persistentVolumeList)
	return &storageClass, err
}
