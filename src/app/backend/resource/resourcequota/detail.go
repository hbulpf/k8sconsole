package resourcequota

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"k8s.io/api/core/v1"
)

// ResourceQuotaDetailList
type ResourceQuotaDetailList struct {
	ListMeta api.ListMeta          `json:"listMeta"`
	Items    []ResourceQuotaDetail `json:"items"`
}


type ResourceQuotaDetail struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`

	// Scopes defines quota scopes
	Scopes []v1.ResourceQuotaScope `json:"scopes,omitempty"`

	// StatusList is a set of (resource name, Used, Hard) tuple.
	StatusList map[v1.ResourceName]ResourceStatus `json:"statusList,omitempty"`
}

// ResourceStatus provides the status of the resource defined by a resource quota.
type ResourceStatus struct {
	Used string `json:"used,omitempty"`
	Hard string `json:"hard,omitempty"`
}

func ToResourceQuotaDetail(rawResourceQuota *v1.ResourceQuota) *ResourceQuotaDetail {
	statusList := make(map[v1.ResourceName]ResourceStatus)

	for key, value := range rawResourceQuota.Status.Hard {
		used := rawResourceQuota.Status.Used[key]
		statusList[key] = ResourceStatus{
			Used: used.String(),
			Hard: value.String(),
		}
	}

	return &ResourceQuotaDetail{
		ObjectMeta: api.NewObjectMeta(rawResourceQuota.ObjectMeta),
		TypeMeta:   api.NewTypeMeta(api.ResourceKindResourceQuota),
		Scopes:     rawResourceQuota.Spec.Scopes,
		StatusList: statusList,
	}
}
