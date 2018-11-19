package persistentvolumeclaim

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

// GetPodPersistentVolumeClaims gets persistentvolumeclaims that are associated with this pod.
func GetPodPersistentVolumeClaims(client kubernetes.Interface, namespace string,
	podName string, dsQuery *dataselect.DataSelectQuery) (*PersistentVolumeClaimList, error) {

	pod, err := client.CoreV1().Pods(namespace).Get(podName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	claimNames := make([]string, 0)
	if pod.Spec.Volumes != nil && len(pod.Spec.Volumes) > 0 {
		for _, v := range pod.Spec.Volumes {
			persistentVolumeClaim := v.PersistentVolumeClaim
			if persistentVolumeClaim != nil {
				claimNames = append(claimNames, persistentVolumeClaim.ClaimName)
			}
		}
	}

	if len(claimNames) > 0 {
		// Get all pvc from namespace
		channels := &common.ResourceChannels{
			PersistentVolumeClaimList: common.GetPersistentVolumeClaimListChannel(client,
				common.NewOneNamespaceQuery(namespace), 1),
		}
		persistentVolumeClaimList := <- channels.PersistentVolumeClaimList.List
		nonCriticalErrors, criticalError := errors.HandleError(err)
		if criticalError != nil {
			return nil, criticalError
		}

		// find pvc in claimNames and persistentVolumeClaimList
		podPersistentVolumeClaims := make([]v1.PersistentVolumeClaim, 0)
		for _, pvc := range persistentVolumeClaimList.Items {
			for _, claimName := range claimNames {
				if strings.Compare(claimName, pvc.Name) == 0 {
					podPersistentVolumeClaims = append(podPersistentVolumeClaims, pvc)
					break
				}
			}
		}

		glog.Infof("Found %d persistentvolumeclaims related to %s pod",
			len(podPersistentVolumeClaims), podName)
		return toPersistentVolumeChaimList(podPersistentVolumeClaims, nonCriticalErrors, dsQuery), nil
	}

	glog.Infof("No persistentvolumeclaims found related to %s pod", podName)
	return &PersistentVolumeClaimList{}, nil
}

type PersistentVolumeClaimCell v1.PersistentVolumeClaim

func (self PersistentVolumeClaimCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
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

func toCells(std []v1.PersistentVolumeClaim) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = PersistentVolumeClaimCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []v1.PersistentVolumeClaim {
	std := make([]v1.PersistentVolumeClaim, len(cells))
	for i := range cells {
		std[i] = v1.PersistentVolumeClaim(cells[i].(PersistentVolumeClaimCell))
	}
	return std
}