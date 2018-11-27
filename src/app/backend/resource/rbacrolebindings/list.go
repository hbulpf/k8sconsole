package rbacrolebindings

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
)

// RbacRoleBindingList contains a list of Roles and ClusterRoles bindings in the cluster.
type RbacRoleBindingList struct {
	ListMeta api.ListMeta `json:"listMeta"`

	// Unordered list of Rbac Role Bindings
	Items []RbacRoleBinding `json:"items"`

	// list of non-critical errors, that occurred during resource retrieval
	Errors []error `json:"errors"`
}

// RbacRoleBinding provides the simplified, combined presentation layer view of Kubernetes' RBAC RoleBindings and ClusterRoleBindings.
// ClusterRoleBindings will be referred to as RoleBindings for the namespace "all namespaces".
type RbacRoleBinding struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`
	Subjects   []rbac.Subject `json:"subjects"`
	RoleRef    rbac.RoleRef   `json:"roleRef"`
	Name       string         `json:"name"`
	Namespace  string         `json:"namespace"`
}

func GetRbacRoleBindingList(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery) (
	*RbacRoleBindingList, error) {
	glog.Info("Getting list rbac role bindings.")
	channels := &common.ResourceChannels{
		RoleBindingList: common.GetRoleBindingListChannel(client, 1),
		ClusterRoleBindingList: common.GetClusterRoleBindingListChannel(client, 1),
	}

	return GetRbacRoleBindingListFromChannels(channels, dsQuery)
}

func GetRbacRoleBindingListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery) (
	*RbacRoleBindingList, error) {

	roleBindings := <- channels.RoleBindingList.List
	err := <- channels.RoleBindingList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	clusterRoleBindings := <- channels.ClusterRoleBindingList.List
	err = <- channels.ClusterRoleBindingList.Error
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	result := toRbacRoleBindingList(roleBindings.Items, clusterRoleBindings.Items, nonCriticalErrors, dsQuery)
	return result, nil
}

func toRbacRoleBindingList(roleBindings []rbac.RoleBinding, clusterRoleBindings []rbac.ClusterRoleBinding,
	nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery) *RbacRoleBindingList {
	result := &RbacRoleBindingList{
		ListMeta: api.ListMeta{TotalItems: len(roleBindings) + len(clusterRoleBindings)},
		Errors: nonCriticalErrors,
	}
	items := make([]RbacRoleBinding, 0)

	for _, item := range roleBindings {
		items = append(items, RbacRoleBinding{
			ObjectMeta: api.NewObjectMeta(item.ObjectMeta),
			TypeMeta: api.NewTypeMeta(api.ResourceKindRbacRoleBinding),
			Subjects: item.Subjects,
			RoleRef: item.RoleRef,
			Name: item.ObjectMeta.Name,
			Namespace: item.ObjectMeta.Namespace,
		})
	}

	for _, item := range clusterRoleBindings {
		items = append(items, RbacRoleBinding{
			ObjectMeta: api.NewObjectMeta(item.ObjectMeta),
			TypeMeta: api.NewTypeMeta(api.ResourceKindRbacClusterRoleBinding),
			Subjects: item.Subjects,
			RoleRef: item.RoleRef,
			Name: item.ObjectMeta.Name,
			Namespace: "",
		})
	}

	roleBindingCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(items), dsQuery)
	result.ListMeta = api.ListMeta{TotalItems: filteredTotal}
	result.Items = fromCells(roleBindingCells)
	return result
}