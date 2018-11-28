package ingress

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	extensions "k8s.io/api/extensions/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// IngressDetail API resource provides mechanisms to inject containers with configuration data while keeping
// containers agnostic of Kubernetes
type IngressDetail struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta api.TypeMeta `json:"typeMeta"`

	Spec extensions.IngressSpec `json:"spec"`

	// Status is the current state of the ingress
	Status extensions.IngressStatus `json:"status"`

	Errors []error `json:"errors"`
}

func GetIngressDetail(client kubernetes.Interface, namespace, name string) (*IngressDetail, error) {
	glog.Infof("Getting details of %s ingress in %s namespace", name, namespace)

	rawIngress, err := client.ExtensionsV1beta1().Ingresses(namespace).Get(name, metaV1.GetOptions{})

	if err != nil {
		return nil, err
	}

	return toIngressDetail(rawIngress), nil
}

func toIngressDetail(ingress *extensions.Ingress) *IngressDetail {
	return &IngressDetail{
		ObjectMeta: api.NewObjectMeta(ingress.ObjectMeta),
		TypeMeta: api.NewTypeMeta(api.ResourceKindIngress),
		Spec: ingress.Spec,
		Status: ingress.Status,
	}
}