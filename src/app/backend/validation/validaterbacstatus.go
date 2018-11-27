package validation

import (
	"fmt"
	auth "k8s.io/api/authentication/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sort"
)

// RbacStatus describe status of RBAC in the cluster
type RbacStatus struct {
	// True when rbac is enable
	Enable bool `json:"enable"`
}

func ValidateRbacStatus(client kubernetes.Interface) (*RbacStatus, error) {
	groupList, err := client.Discovery().ServerGroups()
	if err != nil {
		return nil, fmt.Errorf("Couldn't get available api versions from server: %v", err)
	}

	apiVersions := metaV1.ExtractGroupVersions(groupList)
	return &RbacStatus{
		Enable: contains(apiVersions, auth.SchemeGroupVersion.String()),
	}, nil
}

func contains(arr []string, str string) bool {
	sort.Strings(arr)
	idx := sort.SearchStrings(arr, str)
	return len(arr) > 0 && idx < len(arr) && arr[idx] == str
}
