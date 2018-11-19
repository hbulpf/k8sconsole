package common

import (
	"bytes"
	"k8s.io/api/core/v1"
)

// Endpoint describes an endpoint that is host and a list of available ports for that host
type Endpoint struct {
	// Hostname, either as a domain name or IP address
	Host string `json:"host"`

	// List of ports opened for this endpoint on the hostname
	Ports []ServicePort `json:"ports"`
}

// GetExternalEndpoints returns endpoints that are externally reachable for a service.
func GetExternalEndpoints(service *v1.Service) []Endpoint {
	var externalEndpoints []Endpoint
	if service.Spec.Type == v1.ServiceTypeLoadBalancer {
		for _, ingress := range service.Status.LoadBalancer.Ingress {
			externalEndpoints = append(externalEndpoints, getExternalEndpoint(ingress, service.Spec.Ports))
		}
	}

	for _, ip := range service.Spec.ExternalIPs {
		externalEndpoints = append(externalEndpoints, Endpoint{
			Host: ip,
			Ports: GetServicePorts(service.Spec.Ports),
		})
	}

	return externalEndpoints
}

// GetInternalEndpoint returns internal endpoint name for the given service properties, e.g.,
// "my-service.namespace 80/TCP" or "my-service 53/TCP,53/UDP".
func GetInternalEndpoint(serviceName, namespace string, ports []v1.ServicePort) Endpoint {
	name := serviceName

	if namespace != v1.NamespaceDefault && len(namespace) > 0 && len(serviceName) > 0 {
		bufferName := bytes.NewBufferString(name)
		bufferName.WriteString(".")
		bufferName.WriteString(namespace)
		name = bufferName.String()
	}

	return Endpoint{
		Host: name,
		Ports: GetServicePorts(ports),
	}
}

// Returns external endpoint for the given service properties
func getExternalEndpoint(ingress v1.LoadBalancerIngress, ports []v1.ServicePort) Endpoint {
	var host string
	if ingress.Hostname != "" {
		host = ingress.Hostname
	} else {
		host = ingress.IP
	}
	return Endpoint{
		Host: 	host,
		Ports: 	GetServicePorts(ports),
	}
}