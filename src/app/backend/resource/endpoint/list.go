package endpoint

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type EndpointList struct {
	ListMeta api.ListMeta `json:"listMeta"`
	// List of endpoints
	Endpoints []Endpoint `json:"endpoints"`
}

func GetServiceEndpoints(client kubernetes.Interface, namespace, name string) (*EndpointList, error) {
	endpointList := &EndpointList{
		Endpoints: make([]Endpoint, 0),
		ListMeta: api.ListMeta{TotalItems: 0},
	}

	serviceEndpoints, err := GetEndpoints(client, namespace, name)
	if err != nil {
		return endpointList, err
	}

	endpointList = toEndpointList(serviceEndpoints)
	glog.Infof("Found %d endpoints related to %s service in %s namespace", len(endpointList.Endpoints), name, namespace)
	return endpointList, nil
}

func toEndpointList(endpoints []v1.Endpoints) *EndpointList {
	endpointList := EndpointList{
		Endpoints: make([]Endpoint, 0),
		ListMeta: api.ListMeta{TotalItems: len(endpoints)},
	}

	for _, endpoint := range endpoints {
		for _, subSets := range endpoint.Subsets {
			for _, address := range subSets.Addresses {
				endpointList.Endpoints = append(endpointList.Endpoints, *toEndpoint(address, subSets.Ports, true))
			}
			for _, notReadyAddress := range subSets.NotReadyAddresses {
				endpointList.Endpoints = append(endpointList.Endpoints, *toEndpoint(notReadyAddress, subSets.Ports, false))
			}
		}
	}

	return &endpointList
}
