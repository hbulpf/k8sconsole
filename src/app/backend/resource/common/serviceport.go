package common

import "k8s.io/api/core/v1"

// ServicePort is a pair of port and protocol - service endpoint
type ServicePort struct {
	// Positive port number
	Port int32 `json:"port"`

	// Protocol name, e.g., TCP/UDP
	Protocol v1.Protocol `json:"protocol"`

	// The port on each node on which the service is exposed
	NodePort int32 `json:"nodePort"`
}

// GetServicePorts returns human readable name for the given k8s api service port list.
func GetServicePorts(v1Ports []v1.ServicePort) []ServicePort {
	var ports []ServicePort
	for _, port := range v1Ports {
		ports = append(ports, ServicePort{
			port.Port,
			port.Protocol,
			port.NodePort,
		})
	}
	return ports
}
