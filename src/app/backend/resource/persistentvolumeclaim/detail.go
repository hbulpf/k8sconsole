package persistentvolumeclaim

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PersistentVolumeClaimDetail provides the presentation layer view of Kubernetes Persistent Volume Claim resource.
type PersistentVolumeClaimDetail struct {
	ObjectMeta   api.ObjectMeta                  `json:"objectMeta"`
	TypeMeta     api.TypeMeta                    `json:"typeMeta"`
	Status       v1.PersistentVolumeClaimPhase   `json:"status"`
	Volume       string                          `json:"volume"`
	Capacity     v1.ResourceList                 `json:"capacity"`
	AccessModes  []v1.PersistentVolumeAccessMode `json:"accessModes"`
	StorageClass *string                         `json:"storageClass"`
}

// GetPersistentVolumeClaimDetail returns detailed information about a persistent volume claim
func GetPersistentVolumeClaimDetail(client kubernetes.Interface,
	namespace string, name string) (*PersistentVolumeClaimDetail, error) {
	glog.Infof("Getting details of %s persistent volume claim", name)

	rawPersistentVolumeClaim, err := client.CoreV1().PersistentVolumeClaims(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return getPersistentVolumeClaimDetail(rawPersistentVolumeClaim), nil
}

func getPersistentVolumeClaimDetail(persistentVolumeClaim *v1.PersistentVolumeClaim) *PersistentVolumeClaimDetail {
	return &PersistentVolumeClaimDetail{
		ObjectMeta:   api.NewObjectMeta(persistentVolumeClaim.ObjectMeta),
		TypeMeta:     api.NewTypeMeta(api.ResourceKindPersistentVolumeClaim),
		Status:       persistentVolumeClaim.Status.Phase,
		Volume:       persistentVolumeClaim.Spec.VolumeName,
		Capacity:     persistentVolumeClaim.Status.Capacity,
		AccessModes:  persistentVolumeClaim.Spec.AccessModes,
		StorageClass: persistentVolumeClaim.Spec.StorageClassName,
	}
}