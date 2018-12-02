package horizontalpodautoscaler

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	autoscaling "k8s.io/api/autoscaling/v1"
	"k8s.io/client-go/kubernetes"
)

type HorizontalPodAutoscalerList struct {
	ListMeta api.ListMeta `json:"listMeta"`

	// Unordered list of Horizontal Pod Autoscalers.
	HorizontalPodAutoscalers []HorizontalPodAutoscaler `json:"horizontalpodautoscalers"`

	// List of non-critical errors, that occurred during resource retrieval.
	Errors []error `json:"errors"`
}

// HorizontalPodAutoscaler (aka. Horizontal Pod Autoscaler - hpa)
type HorizontalPodAutoscaler struct {
	ObjectMeta                      api.ObjectMeta `json:"objectMeta"`
	TypeMeta                        api.TypeMeta   `json:"typeMeta"`
	ScaleTargetRef                  ScaleTargetRef `json:"scaleTargetRef"`
	MinReplicas                     *int32         `json:"minReplicas"`
	MaxReplicas                     int32          `json:"maxReplicas"`
	CurrentCPUUtilizationPercentage *int32         `json:"currentCPUUtilizationPercentage"`
	TargetCPUUtilizationPercentage  *int32         `json:"targetCPUUtilizationPercentage"`
}

func GetHorizontalPodAutoscalerList(client kubernetes.Interface, nsQuery *common.NamespaceQuery) (*HorizontalPodAutoscalerList, error) {
	channel := common.GetHorizontalPodAutoscalerListChannel(client, nsQuery, 1)
	hpaList := <-channel.List
	err := <-channel.Error

	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return toHorizontalPodAutoscalerList(hpaList.Items, nonCriticalErrors), nil
}

func GetHorizontalPodAutoscalerListForResource(client kubernetes.Interface, namespace, kind, name string) (*HorizontalPodAutoscalerList, error) {
	nsQuery := common.NewOneNamespaceQuery(namespace)
	channel := common.GetHorizontalPodAutoscalerListChannel(client, nsQuery, 1)
	hpaList := <-channel.List
	err := <-channel.Error

	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	filteredHpaList := make([]autoscaling.HorizontalPodAutoscaler, 0)
	for _, hpa := range hpaList.Items {
		if hpa.Spec.ScaleTargetRef.Kind == kind && hpa.Spec.ScaleTargetRef.Name == name {
			filteredHpaList = append(filteredHpaList, hpa)
		}
	}

	return toHorizontalPodAutoscalerList(filteredHpaList, nonCriticalErrors), nil
}

func toHorizontalPodAutoscalerList(hpas []autoscaling.HorizontalPodAutoscaler, nonCriticalErrors []error) *HorizontalPodAutoscalerList {
	hpaList := &HorizontalPodAutoscalerList{
		HorizontalPodAutoscalers: make([]HorizontalPodAutoscaler, 0),
		ListMeta:                 api.ListMeta{TotalItems: len(hpas)},
		Errors:                   nonCriticalErrors,
	}

	for _, hpa := range hpas {
		horizontalPodAutoscaler := toHorizontalPodAutoScaler(&hpa)
		hpaList.HorizontalPodAutoscalers = append(hpaList.HorizontalPodAutoscalers, horizontalPodAutoscaler)
	}
	return hpaList
}

func toHorizontalPodAutoScaler(hpa *autoscaling.HorizontalPodAutoscaler) HorizontalPodAutoscaler {
	return HorizontalPodAutoscaler{
		ObjectMeta: api.NewObjectMeta(hpa.ObjectMeta),
		TypeMeta:   api.NewTypeMeta(api.ResourceKindHorizontalPodAutoscaler),
		ScaleTargetRef: ScaleTargetRef{
			Kind: hpa.Spec.ScaleTargetRef.Kind,
			Name: hpa.Spec.ScaleTargetRef.Name,
		},
		MinReplicas:                     hpa.Spec.MinReplicas,
		MaxReplicas:                     hpa.Spec.MaxReplicas,
		CurrentCPUUtilizationPercentage: hpa.Status.CurrentCPUUtilizationPercentage,
		TargetCPUUtilizationPercentage:  hpa.Spec.TargetCPUUtilizationPercentage,
	}

}