package secret

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// SecretSpec is a common interface for the specification of different secrets.
type SecretSpec interface {
	GetName() string
	GetType() v1.SecretType
	GetNamespace() string
	GetData() map[string][]byte
}

// CreateSecret creates a single secret using the cluster API client
func CreateSecret(client kubernetes.Interface, spec SecretSpec) (*Secret, error) {
	namespace := spec.GetNamespace()
	secret := &v1.Secret{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      spec.GetName(),
			Namespace: namespace,
		},
		Type: spec.GetType(),
		Data: spec.GetData(),
	}
	_, err := client.CoreV1().Secrets(namespace).Create(secret)
	return toSecret(secret), err
}

// ImagePullSecretSpec is a specification of an image pull secret implements SecretSpec
type ImagePullSecretSpec struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`

	// The value of the .dockercfg property. It must be Base64 encoded.
	Data []byte `json:"data"`
}

// GetName returns the name of the ImagePullSecret
func (spec *ImagePullSecretSpec) GetName() string {
	return spec.Name
}

// GetType returns the type of the ImagePullSecret, which is always api.SecretTypeDockercfg
func (spec *ImagePullSecretSpec) GetType() v1.SecretType {
	return v1.SecretTypeDockercfg
}

// GetNamespace returns the namespace of the ImagePullSecret
func (spec *ImagePullSecretSpec) GetNamespace() string {
	return spec.Namespace
}

// GetData returns the data the secret carries, it is a single key-value pair
func (spec *ImagePullSecretSpec) GetData() map[string][]byte {
	return map[string][]byte{v1.DockerConfigKey: spec.Data}
}

// Secret is a single secret returned to the frontend.
type Secret struct {
	ObjectMeta api.ObjectMeta `json:"objectMeta"`
	TypeMeta   api.TypeMeta   `json:"typeMeta"`
	Type       v1.SecretType  `json:"type"`
}

// SecretsList is a response structure for a queried secrets list.
type SecretList struct {
	api.ListMeta `json:"listMeta"`

	// Unordered list of Items.
	Secrets []Secret `json:"secrets"`

	// List of non-critical errors, that occurred during resource retrieval.
	Errors []error `json:"errors"`
}

func GetSecretList(client kubernetes.Interface, namespace *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*SecretList, error) {

	glog.Infof("Getting list of secrets in %s namespace\n", namespace)
	secretList, err := client.CoreV1().Secrets(namespace.ToRequestParam()).List(api.ListEverything)

	nonCriticalErrors, criticalError := errors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return toSecretList(secretList.Items, nonCriticalErrors, dsQuery), nil
}

func GetSecretListFromChannels(channels *common.ResourceChannels, dsQuery *dataselect.DataSelectQuery) (
	*SecretList, error) {
  secretList := <- channels.SecretList.List
  err := <- channels.SecretList.Error

  nonCriticalErrors, criticalError := errors.HandleError(err)
  if criticalError != nil {
  	return nil, criticalError
	}

  return toSecretList(secretList.Items, nonCriticalErrors, dsQuery), nil
}

func toSecret(secret *v1.Secret) *Secret {
	return &Secret{
		ObjectMeta: api.NewObjectMeta(secret.ObjectMeta),
		TypeMeta:   api.NewTypeMeta(api.ResourceKindSecret),
		Type:       secret.Type,
	}
}

func toSecretList(secrets []v1.Secret, nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery) *SecretList {
	result := &SecretList{
		Secrets:  make([]Secret, 0),
		ListMeta: api.ListMeta{TotalItems: len(secrets)},
		Errors:   nonCriticalErrors,
	}

	secretCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(secrets), dsQuery)
	secrets = fromCells(secretCells)
	result.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, secret := range secrets {
		result.Secrets = append(result.Secrets, *toSecret(&secret))
	}

	return result
}