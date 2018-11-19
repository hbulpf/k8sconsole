package api

import (
	"github.com/emicklei/go-restful"
	authApi "github.com/wzt3309/k8sconsole/src/app/backend/auth/api"
	"k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// ClientManager is responsible for initializing and creating clients to communicate with
// k8s apiserver.
type ClientManager interface {
	Client(*restful.Request) (kubernetes.Interface, error)
	InsecureClient() kubernetes.Interface
	CanI(*restful.Request, *v1.SelfSubjectAccessReview) bool
	Config(*restful.Request) (*rest.Config, error)
	ClientCmdConfig(*restful.Request) (clientcmd.ClientConfig, error)
	CSRFKey() string
	HasAccess(api.AuthInfo) error
	VerberClient(*restful.Request) (ResourceVerber, error)
	SetTokenManager(manager authApi.TokenManager)
}

// ResourceVerber is responsible for performing generic CRUD operations on all supported resources.
type ResourceVerber interface {
	Put(kind string, namespaceSet bool, namespace string, name string, object *runtime.Unknown) error
	Get(kind string, namespaceSet bool, namespace string, name string) (runtime.Object, error)
	Delete(kind string, namespaceSet bool, namespace string, name string) error
}

// CanIResponse represents a response that contains the result of checking whether or not user is allowed
// to access the given resource.
type CanIResponse struct {
	Allowed bool `json:"allowed"`
}
