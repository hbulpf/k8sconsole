package configmap

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ConfigMapList contains a list of Config Maps in the cluster.
type ConfigMapList struct {
	ListMeta api.ListMeta `json:"listMeta"`

	// Unordered list of Config Maps
	Items []ConfigMap `json:"items"`

	// List of non-critical errors, that occurred during resource retrieval.
	Errors []error `json:"errors"`
}

// ConfigMap API resource provides mechanisms to inject containers with configuration data while keeping
// containers agnostic of Kubernetes
type ConfigMap struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`
}

// GetConfigMapList returns a list of all ConfigMaps in the cluster.
func GetConfigMapList(client kubernetes.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*ConfigMapList, error) {
	glog.Infof("Getting list config maps in the namespace %s", nsQuery.ToRequestParam())
	channels := &common.ResourceChannels{
		ConfigMapList: common.GetConfigMapListChannel(client, nsQuery, 1),
	}
	return GetConfigMapListFromChannels(channels, dsQuery)
}

func GetConfigMapListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery) (
	*ConfigMapList, error) {
	configMaps := <- channels.ConfigMapList.List
	err := <- channels.ConfigMapList.Error

	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return toConfigMapList(configMaps.Items, nonCriticalErrors, dsQuery), nil
}

func toConfigMapList(configMaps []v1.ConfigMap, nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery) *ConfigMapList {
	result := &ConfigMapList{
		Items: make([]ConfigMap, 0),
		ListMeta: api.ListMeta{TotalItems: len(configMaps)},
		Errors: nonCriticalErrors,
	}

	configMapCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(configMaps), dsQuery)
	configMaps = fromCells(configMapCells)
	result.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, item := range configMaps {
		result.Items = append(result.Items, toConfigMap(item.ObjectMeta))
	}

	return result
}

func toConfigMap(meta metaV1.ObjectMeta) ConfigMap {
	return ConfigMap{
		ObjectMeta: api.NewObjectMeta(meta),
		TypeMeta: api.NewTypeMeta(api.ResourceKindConfigMap),
	}
}
