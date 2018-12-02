package horizontalpodautoscaler

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	autoscaling "k8s.io/api/autoscaling/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// HorizontalPodAutoscalerDetail provides the presentation layer view of Kubernetes Horizontal Pod Autoscaler resource.
// close mapping of the autoscaling.HorizontalPodAutoscaler type with part of the *Spec and *Detail childs
type HorizontalPodAutoscalerDetail struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`

	ScaleTargetRef ScaleTargetRef `json:"scaleTargetRef"`

	MinReplicas *int32 `json:"minReplicas"`
	MaxReplicas int32  `json:"maxReplicas"`

	CurrentCPUUtilizationPercentage *int32 `json:"currentCPUUtilizationPercentage"`
	TargetCPUUtilizationPercentage  *int32 `json:"targetCPUUtilizationPercentage"`

	CurrentReplicas int32 `json:"currentReplicas"`
	DesiredReplicas int32 `json:"desiredReplicas"`

	LastScaleTime *v1.Time `json:"lastScaleTime"`
}

// GetHorizontalPodAutoscalerDetail returns detailed information about a horizontal pod autoscaler
func GetHorizontalPodAutoscalerDetail(client kubernetes.Interface, namespace string, name string) (*HorizontalPodAutoscalerDetail, error) {
	glog.Infof("Getting details of %s horizontal pod autoscaler", name)

	rawHorizontalPodAutoscaler, err := client.AutoscalingV1().HorizontalPodAutoscalers(namespace).Get(name, v1.GetOptions{})

	if err != nil {
		return nil, err
	}

	return getHorizontalPodAutoscalerDetail(rawHorizontalPodAutoscaler), nil
}

func getHorizontalPodAutoscalerDetail(horizontalPodAutoscaler *autoscaling.HorizontalPodAutoscaler) *HorizontalPodAutoscalerDetail {

	return &HorizontalPodAutoscalerDetail{
		ObjectMeta: api.NewObjectMeta(horizontalPodAutoscaler.ObjectMeta),
		TypeMeta:   api.NewTypeMeta(api.ResourceKindHorizontalPodAutoscaler),

		ScaleTargetRef: ScaleTargetRef{
			Kind: horizontalPodAutoscaler.Spec.ScaleTargetRef.Kind,
			Name: horizontalPodAutoscaler.Spec.ScaleTargetRef.Name,
		},

		MinReplicas:                     horizontalPodAutoscaler.Spec.MinReplicas,
		MaxReplicas:                     horizontalPodAutoscaler.Spec.MaxReplicas,
		CurrentCPUUtilizationPercentage: horizontalPodAutoscaler.Status.CurrentCPUUtilizationPercentage,
		TargetCPUUtilizationPercentage:  horizontalPodAutoscaler.Spec.TargetCPUUtilizationPercentage,

		CurrentReplicas: horizontalPodAutoscaler.Status.CurrentReplicas,
		DesiredReplicas: horizontalPodAutoscaler.Status.DesiredReplicas,

		LastScaleTime: horizontalPodAutoscaler.Status.LastScaleTime,
	}
}