package ingress

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	extensions "k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"
)

type Ingress struct {
	api.ObjectMeta `json:"objectMeta"`
	api.TypeMeta `json:"typeMeta"`

	// External endpoints of this ingress
	Endpoints []common.Endpoint `json:"endpoints"`
}

type IngressList struct {
	api.ListMeta `json:"listMeta"`

	// Unordered list of ingress
	Items []Ingress `json:"items"`

	// List of non-critical errors, that occurred during resource retrieval
	Errors []error `json:"errors"`
}

// GetIngressList returns all ingresses in the given namespace.
func GetIngressList(client kubernetes.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*IngressList, error) {
	ingressList, err := client.ExtensionsV1beta1().Ingresses(nsQuery.ToRequestParam()).List(api.ListEverything)

	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return toIngressList(ingressList.Items, nonCriticalErrors, dsQuery), nil
}

// GetIngressListFromChannels - return all ingresses in the given namespace.
func GetIngressListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery) (
	*IngressList, error) {
	ingresses := <- channels.IngressList.List
	err := <- channels.IngressList.Error

	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return toIngressList(ingresses.Items, nonCriticalErrors, dsQuery), nil
}

func toIngressList(ingresses []extensions.Ingress, nonCriticalErrors []error,
	dsQuery *dataselect.DataSelectQuery) *IngressList {
	newIngressList := &IngressList{
		ListMeta: api.ListMeta{TotalItems: len(ingresses)},
		Items: make([]Ingress, 0),
		Errors: nonCriticalErrors,
	}

	ingressCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(ingresses), dsQuery)
	newIngressList.ListMeta = api.ListMeta{TotalItems: filteredTotal}
	ingresses = fromCells(ingressCells)

	for _, ingress := range ingresses {
		newIngressList.Items = append(newIngressList.Items, toIngress(ingress))
	}

	return newIngressList
}

func toIngress(ingress extensions.Ingress) Ingress {
	return Ingress{
		ObjectMeta: api.NewObjectMeta(ingress.ObjectMeta),
		TypeMeta: api.NewTypeMeta(api.ResourceKindIngress),
		Endpoints: getEndpoints(ingress),
	}
}

func getEndpoints(ingress extensions.Ingress) []common.Endpoint {
	endpoints := make([]common.Endpoint, 0)
	if len(ingress.Status.LoadBalancer.Ingress) > 0 {
		for _, status := range ingress.Status.LoadBalancer.Ingress {
			endpoint := common.Endpoint{Host: status.IP}
			endpoints = append(endpoints, endpoint)
		}
	}
	return endpoints
}
