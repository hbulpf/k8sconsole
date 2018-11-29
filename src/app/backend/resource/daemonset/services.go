package daemonset

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/service"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetDaemonSetServices returns list of services that are related to daemon set targeted by given
// name.
func GetDaemonSetServices(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery,
	namespace, name string) (*service.ServiceList, error) {

	daemonSet, err := client.AppsV1beta2().DaemonSets(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	channels := &common.ResourceChannels{
		ServiceList: common.GetServiceListChannel(client, common.NewOneNamespaceQuery(namespace), 1),
	}

	services := <-channels.ServiceList.List
	err = <-channels.ServiceList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	matchingServices := common.FilterNamespacedServicesBySelector(services.Items, namespace,
		daemonSet.Spec.Selector.MatchLabels)
	return service.ToServiceList(matchingServices, nonCriticalErrors, dsQuery), nil
}