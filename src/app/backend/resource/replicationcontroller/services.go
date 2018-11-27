package replicationcontroller

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/service"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetReplicationControllerServices(client kubernetes.Interface, dsQuery *dataselect.DataSelectQuery,
	namespace, rcName string) (*service.ServiceList, error) {

	replicationController, err := client.CoreV1().ReplicationControllers(namespace).Get(rcName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	channels := &common.ResourceChannels{
		ServiceList: common.GetServiceListChannel(client, common.NewOneNamespaceQuery(namespace), 1),
	}

	services := <- channels.ServiceList.List
	err = <- channels.ServiceList.Error
	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	matchingServices := common.FilterNamespacedServicesBySelector(services.Items,
		namespace, replicationController.Spec.Selector)
	return service.ToServiceList(matchingServices, nonCriticalErrors, dsQuery), nil
}

