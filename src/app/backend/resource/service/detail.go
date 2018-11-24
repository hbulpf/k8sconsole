package service

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/endpoint"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/pod"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ServiceDetail struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`

	// InternalEndpoint of all Kubernetes services that have the same label selector as connected Replication
	// Controller. Endpoints is DNS name merged with ports.
	InternalEndpoint common.Endpoint `json:"internalEndpoint"`

	// ExternalEndpoints of all Kubernetes services that have the same label selector as connected Replication
	// Controller. Endpoints is external IP address name merged with ports.
	ExternalEndpoints []common.Endpoint `json:"externalEndpoints"`

	// List of Endpoint obj. that are endpoints of this Service.
	EndpointList endpoint.EndpointList `json:"endpointList"`

	// Label selector of the service.
	Selector map[string]string `json:"selector"`

	// Type determines how the service will be exposed.  Valid options: ClusterIP, NodePort, LoadBalancer
	Type v1.ServiceType `json:"type"`

	// ClusterIP is usually assigned by the master. Valid values are None, empty string (""), or
	// a valid IP address. None can be specified for headless services when proxying is not required
	ClusterIP string `json:"clusterIP"`

	// List of events related to this Service
	EventList common.EventList `json:"eventList"`

	// PodList represents list of pods targeted by same label selector as this service.
	PodList pod.PodList `json:"podList"`

	// Show the value of the SessionAffinity of the Service.
	SessionAffinity v1.ServiceAffinity `json:"sessionAffinity"`

	// List of non-critical errors, that occurred during resource retrieval.
	Errors []error `json:"errors"`
}

func GetServiceDetail(client kubernetes.Interface, namespace, name string,
	dsQuery *dataselect.DataSelectQuery) (*ServiceDetail, error) {

		glog.Infof("Getting details of %s service in %s namespace", name, namespace)
		serviceData, err := client.CoreV1().Services(namespace).Get(name, metaV1.GetOptions{})
		if err != nil {
			return nil, err
		}

		endpointList, err := endpoint.GetServiceEndpoints(client, namespace, name)
		nonCriticalErrors, criticalError := errors.HandleError(err)
		if criticalError != nil {
			return nil, criticalError
		}

		podList, err := GetServicePods(client, namespace, name, dsQuery)
		nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
		if criticalError != nil {
			return nil, criticalError
		}

		eventList, err := GetServiceEvents(client, dsQuery, namespace, name)
		nonCriticalErrors, criticalError = errors.AppendError(err, nonCriticalErrors)
		if criticalError != nil {
			return nil, criticalError
		}

		service := ToServiceDetail(serviceData, *eventList, *podList, *endpointList, nonCriticalErrors)
		return &service, nil
}

// ToServiceDetail returns api service object based on kubernetes service object
func ToServiceDetail(service *v1.Service, events common.EventList, pods pod.PodList,
	endpointList endpoint.EndpointList, nonCriticalErrors []error) ServiceDetail {
	return ServiceDetail{
		ObjectMeta: 				api.NewObjectMeta(service.ObjectMeta),
		TypeMeta: 					api.NewTypeMeta(api.ResourceKindService),
		InternalEndpoint: 	common.GetInternalEndpoint(service.Name, service.Namespace, service.Spec.Ports),
		ExternalEndpoints: 	common.GetExternalEndpoints(service),
		EndpointList: 			endpointList,
		Selector: 					service.Spec.Selector,
		ClusterIP: 					service.Spec.ClusterIP,
		Type: 							service.Spec.Type,
		EventList: 					events,
		PodList: 						pods,
		SessionAffinity: 		service.Spec.SessionAffinity,
		Errors: 						nonCriticalErrors,
	}
}