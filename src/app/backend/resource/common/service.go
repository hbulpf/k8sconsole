package common

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"k8s.io/api/core/v1"
)

// FilterNamespacedServicesBySelector returns services targeted by given resource selector in
// given namespace.
func FilterNamespacedServicesBySelector(services []v1.Service, namespace string,
	resourceSelector map[string]string) []v1.Service {
	var matchingServices []v1.Service
	for _, service := range services {
		if service.ObjectMeta.Namespace == namespace &&
			api.IsSelectorMatching(service.Spec.Selector, resourceSelector){
			matchingServices = append(matchingServices, service)
		}
	}

	return matchingServices
}
