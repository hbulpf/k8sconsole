package validation

import (
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
)

// ProtocolValiditySpec is a specification of protocol validation request.
type ProtocolValiditySpec struct {
	// Protocol type
	Protocol v1.Protocol `json:"protocol"`

	// Service type. LoadBalancer(true)/NodePort(false).
	IsExternal bool `json:"isExternal"`
}

// ProtocolValidity describes validity of the protocol.
type ProtocolValidity struct {
	// True when the selected protocol is valid for selected service type.
	Valid bool `json:"valid"`
}

// ValidateProtocol validates protocol based on whether created service is set to NodePort or NodeBalancer type.
func ValidateProtocol(spec *ProtocolValiditySpec) *ProtocolValidity {
	glog.Infof("Validating %s protocol for service with external set to %v", spec.Protocol, spec.IsExternal)

	isValid := true
	if spec.Protocol == v1.ProtocolUDP && spec.IsExternal {
		isValid = false
	}

	glog.Infof("Validation result for %s protocol is %v", spec.Protocol, isValid)
	return &ProtocolValidity{Valid: isValid}
}