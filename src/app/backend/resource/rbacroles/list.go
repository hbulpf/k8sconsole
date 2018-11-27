package rbacroles

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	rbac "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// RbacRoleList contains a list of roles and cluster roles in the cluster
type RbacRoleList struct {
	ListMeta api.ListMeta `json:"listMeta"`

	// Unordered list of rbac role
	Items []RbacRole `json:"items"`

	// list of non-critical errors, that occurred during resource retrieval
	Errors []error `json:"errors"`
}

type RbacRole struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta api.TypeMeta `json:"typeMeta"`
}

func GetRbacRoleList(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery) (*RbacRoleList, error) {
	glog.Info("Getting list of RBAC roles.")
	channels := &common.ResourceChannels{
		RoleList: common.GetRoleListChannel(client, 1),
		ClusterRoleList: common.GetClusterRoleListChannel(client, 1),
	}

	return GetRbacRoleListFromChannels(channels, dsQuery)
}

func GetRbacRoleListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery) (
	*RbacRoleList, error) {
	roles := <- channels.RoleList.List
	err := <- channels.RoleList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	clusterRoles := <- channels.ClusterRoleList.List
	err = <- channels.ClusterRoleList.Error
	nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
	if criticalError != nil {
		return nil, criticalError
	}

	result := toRbacRoleList(roles.Items, clusterRoles.Items, nonCriticalErrors, dsQuery)
	return result, nil
}

func toRbacRoleList(roles []rbac.Role, clusterRoles []rbac.ClusterRole, nonCriticalErrors []error,
	dsQuery *dataselect.DataSelectQuery) *RbacRoleList {
	result := &RbacRoleList{
		ListMeta: api.ListMeta{TotalItems: len(roles) + len(clusterRoles)},
		Errors: nonCriticalErrors,
	}

	items := make([]RbacRole, 0)
	for _, role := range roles {
		items = append(items, toRbacRole(role.ObjectMeta, api.ResourceKindRbacRole))
	}

	for _, clusterRole := range clusterRoles {
		items = append(items, toRbacRole(clusterRole.ObjectMeta, api.ResourceKindRbacClusterRole))
	}

	roleCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(items), dsQuery)
	result.ListMeta = api.ListMeta{TotalItems: filteredTotal}
	result.Items = fromCells(roleCells)

	return result
}

func toRbacRole(meta v1.ObjectMeta, kind api.ResourceKind) RbacRole {
	return RbacRole{
		ObjectMeta: api.NewObjectMeta(meta),
		TypeMeta: api.NewTypeMeta(kind),
	}
}
