package validation

import (
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// AppNameValiditySpec is a specification for application name validation request
type AppNameValiditySpec struct {
	Name string `json:"name"`
	Namespace string `json:"namespace"`
}

// AppNameValidity describes validity of the application name
type AppNameValidity struct {
	// True when the application name is valid
	Valid bool `json:"valid"`
}

func ValidateAppName(spec *AppNameValiditySpec, client kubernetes.Interface) (*AppNameValidity, error) {
	glog.Infof("Validating %s application name in %s namespace", spec.Name, spec.Namespace)

	isValidDeployment := false
	isValidService := false

	_, err := client.AppsV1beta2().Deployments(spec.Namespace).Get(spec.Name, metaV1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) || errors.IsForbidden(err) {
			isValidDeployment = true
		} else {
			return nil, err
		}
	}

	_, err = client.CoreV1().Services(spec.Namespace).Get(spec.Name, metaV1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) || errors.IsForbidden(err) {
			isValidService = true
		} else {
			return nil, err
		}
	}

	isValid := isValidDeployment && isValidService

	glog.Infof("Validation result for %s application name in %s namespace is %t",
		spec.Name, spec.Namespace, isValid)

	return &AppNameValidity{Valid: isValid}, nil
}
