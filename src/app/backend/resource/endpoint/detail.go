package endpoint

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type Endpoint struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta api.TypeMeta `json:"typeMeta"`

	// Hostname, either as a domain name or IP address
	Host string `json:"host"`

	// Name of the node the endpoint is located
	NodeName *string `json:"nodeName"`

	// Status of the endpoint
	Ready bool `json:"ready"`

	// Array of endpoint ports
	Ports []v1.EndpointPort `json:"ports"`
}

func GetEndpoints(client kubernetes.Interface, namespace, name string) ([]v1.Endpoints, error) {
	fieldSelector, err := fields.ParseSelector("metadata.name" + "=" + name)
	if err != nil {
		return nil, err
	}

	channels := &common.ResourceChannels{
		EndpointList: common.GetEndpointListChannelWithOptions(client,
			common.NewOneNamespaceQuery(namespace),
			metaV1.ListOptions{
				LabelSelector: labels.Everything().String(),
				FieldSelector: fieldSelector.String(),
			},
			1),
	}

	endpointList := <- channels.EndpointList.List
	if err := <- channels.EndpointList.Error; err != nil {
		return nil, err
	}

	return endpointList.Items, nil
}

// toEndpoint converts endpoint api Endpoint to Endpoint model object.
func toEndpoint(address v1.EndpointAddress, ports []v1.EndpointPort, ready bool) *Endpoint {
	return &Endpoint{
		TypeMeta: api.NewTypeMeta(api.ResourceKindEndpoint),
		Host: address.IP,
		Ports: ports,
		Ready: ready,
		NodeName: address.NodeName,
	}
}
