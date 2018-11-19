package persistentvolume

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

func GetStorageClassPersistentVolumes(client kubernetes.Interface, storageClassName string,
	dsQuery *dataselect.DataSelectQuery) (*PersistentVolumeList, error) {

	storageClass, err := client.StorageV1().StorageClasses().Get(storageClassName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	channels := &common.ResourceChannels{
		PersistentVolumeList: common.GetPersistentVolumeListChannel(client, 1),
	}

	persistentVolumeList := <- channels.PersistentVolumeList.List
	err = <- channels.PersistentVolumeList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	storagePersistentVolumes := make([]v1.PersistentVolume, 0)
	for _, pv := range persistentVolumeList.Items {
		if strings.Compare(pv.Spec.StorageClassName, storageClass.Name) == 0 {
			storagePersistentVolumes = append(storagePersistentVolumes, pv)
		}
	}

	glog.Infof("Found %d persistentvolumes related to %s storageclass",
		len(storagePersistentVolumes), storageClassName)

	return toPersistentVolumeList(storagePersistentVolumes, nonCriticalErrors, dsQuery), nil
}

// getPersistentVolumeClaim returns Persistent Volume claim using "namespace/claim" format.
func getPersistentVolumeClaim(pv *v1.PersistentVolume) string {
	var claim string

	if pv.Spec.ClaimRef != nil {
		claim = pv.Spec.ClaimRef.Namespace + "/" + pv.Spec.ClaimRef.Name
	}
	return claim
}

type PersistentVolumeCell v1.PersistentVolume

func (self PersistentVolumeCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Name)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(self.ObjectMeta.CreationTimestamp.Time)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Namespace)
	default:
		return nil
	}
}

func toCells(std []v1.PersistentVolume) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = PersistentVolumeCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []v1.PersistentVolume {
	std := make([]v1.PersistentVolume, len(cells))
	for i := range cells {
		std[i] = v1.PersistentVolume(cells[i].(PersistentVolumeCell))
	}
	return std
}
