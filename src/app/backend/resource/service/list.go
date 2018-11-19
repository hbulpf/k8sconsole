package service

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Service struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta api.TypeMeta `json:"typeMeta"`

	// InternalEndpoint is DNS name merged with ports
	InternalEndpoint common.Endpoint `json:"internalEndpoint"`

	// ExternalEndpoints are external IP address name merged with ports.
	ExternalEndpoints []common.Endpoint `json:"externalEndpoints"`

	// Label selector of the service
	Selector map[string]string `json:"selector"`

	// Type determines how the service will exposed. Valid options: ClusterIP, NodePort, LoadBalancer
	Type v1.ServiceType `json:"type"`

	// ClusterIP is usually assigned by the master.
	// Valid values:
	// - None (can be specified for headless services when proxying is not required)
	// - empty string ("")
	// - valid IP address
	ClusterIP string `json:"clusterIP"`
}

// ServiceList contains a list of services in the cluster
type ServiceList struct {
	ListMeta api.ListMeta `json:"listMeta"`

	// Unordered list of services.
	Services []Service `json:"services"`

	Errors []error `json:"errors"`
}

// GetServiceList returns a list of all services in the cluster.
func GetServiceList(client kubernetes.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*ServiceList, error) {
	glog.Info("Getting list of all services in the cluster")

	channels := &common.ResourceChannels{
		ServiceList: common.GetServiceListChannel(client, nsQuery, 1),
	}

	return GetServiceListFromChannels(channels, dsQuery)
}

// GetServiceListFromChannels returns a list of all services in the cluster.
func GetServiceListFromChannels(channels *common.ResourceChannels,
	dsQuery *dataselect.DataSelectQuery) (*ServiceList, error) {

	serviceList := <- channels.ServiceList.List
	err := <- channels.ServiceList.Error

	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return ToServiceList(serviceList.Items, nonCriticalErrors, dsQuery), nil
}

// ToService returns api service object based on kubernetes service object
func ToService(service *v1.Service) Service {
	return Service{
		ObjectMeta: api.NewObjectMeta(service.ObjectMeta),
		TypeMeta: api.NewTypeMeta(api.ResourceKindService),
		InternalEndpoint: common.GetInternalEndpoint(service.Name, service.Namespace, service.Spec.Ports),
		ExternalEndpoints: common.GetExternalEndpoints(service),
		Selector: service.Spec.Selector,
		Type: service.Spec.Type,
		ClusterIP: service.Spec.ClusterIP,
	}
}

// ToServiceList returns paginated service list based on given service array and pagination query.
func ToServiceList(services []v1.Service, nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery) *ServiceList {
	serviceList := &ServiceList{
		Services: make([]Service, 0),
		ListMeta: api.ListMeta{TotalItems: len(services)},
		Errors: nonCriticalErrors,
	}

	serviceCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(services), dsQuery)
	services = fromCells(serviceCells)
	serviceList.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, service := range services {
		serviceList.Services = append(serviceList.Services, ToService(&service))
	}

	return serviceList
}

