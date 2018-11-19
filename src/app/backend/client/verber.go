package client

import (
	"fmt"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	clientApi "github.com/wzt3309/k8sconsole/src/app/backend/client/api"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

// resourceVerber is a struct responsible for doing CRUD operations on resources
type resourceVerber struct {
	client           RESTClient
	extensionsClient RESTClient
}

func (verber *resourceVerber) getRESTClientByType(clientType api.ClientType) RESTClient {
	switch clientType {
	case api.ClientTypeExtensionClient:
		return verber.extensionsClient
	default:
		return verber.client
	}
}

// RESTClient is an interface for REST operations used in the file
type RESTClient interface {
	Put() *rest.Request
	Get() *rest.Request
	Delete() *rest.Request
}

// Put puts new resource version of the given kind in the given namespace with the given name.
func (verber *resourceVerber) Put(kind string, namespaceSet bool, namespace string, name string,
	object *runtime.Unknown) error {
	resourceSpec, ok := api.KindToAPIMapping[kind]
	if !ok {
		return fmt.Errorf("Unknown resource kind: %s", kind)
	}

	if namespaceSet != resourceSpec.Namespaced {
		if namespaceSet {
			return fmt.Errorf("Set namespace for not-namespaced resource kind: %s", kind)
		} else {
			return fmt.Errorf("Set no-namespaced for namespaced resource kind: %s", kind)
		}
	}

	client := verber.getRESTClientByType(resourceSpec.ClientType)

	req := client.Put().
		Resource(resourceSpec.Resource).
		Name(name).
		SetHeader("Content-Type", "application/json").
		Body([]byte(object.Raw))

	if resourceSpec.Namespaced {
		req.Namespace(namespace)
	}

	return req.Do().Error()
}

// Get gets the resource of the given kind in the given namespace with the given name.
func (verber *resourceVerber) Get(kind string, namespaceSet bool, namespace string, name string) (
	runtime.Object, error) {
	resourceSpec, ok := api.KindToAPIMapping[kind]
	if !ok {
		return nil, fmt.Errorf("Unknown resource kind: %s", kind)
	}

	if namespaceSet != resourceSpec.Namespaced {
		if namespaceSet {
			return nil, fmt.Errorf("Set namespace for not-namespaced resource kind: %s", kind)
		} else {
			return nil, fmt.Errorf("Set no-namespaced for namespaced resource kind: %s", kind)
		}
	}

	client := verber.getRESTClientByType(resourceSpec.ClientType)
	result := &runtime.Unknown{}
	req := client.Get().
		Resource(resourceSpec.Resource).
		Name(name).
		SetHeader("Accept", "application/json")

	if resourceSpec.Namespaced {
		req.Namespace(namespace)
	}

	err := req.Do().Into(result)
	return result, err
}

// Delete deletes the resource of the given kind in the given namespace with the given name.
func (verber *resourceVerber) Delete(kind string, namespaceSet bool, namespace string, name string) error {
	resourceSpec, ok := api.KindToAPIMapping[kind]
	if !ok {
		return fmt.Errorf("Unknown resource kind: %s", kind)
	}

	if namespaceSet != resourceSpec.Namespaced {
		if namespaceSet {
			return fmt.Errorf("Set namespace for not-namespaced resource kind: %s", kind)
		} else {
			return fmt.Errorf("Set no namespace for namespaced resource kind: %s", kind)
		}
	}

	client := verber.getRESTClientByType(resourceSpec.ClientType)

	// Do cascade delete by default, as this is what users typically expect.
	defaultPropagationPolicy := v1.DeletePropagationForeground
	defaultDeleteOptions := &v1.DeleteOptions{
		PropagationPolicy: &defaultPropagationPolicy,
	}

	req := client.Delete().Resource(resourceSpec.Resource).Name(name).Body(defaultDeleteOptions)

	if resourceSpec.Namespaced {
		req.Namespace(namespace)
	}

	return req.Do().Error()
}

func NewResourceVerber(client, extensionsClient RESTClient) clientApi.ResourceVerber {
	return &resourceVerber{
		client, extensionsClient,
	}
}
