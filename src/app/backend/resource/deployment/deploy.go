package deployment

import (
	"fmt"
	"github.com/golang/glog"
	"io"
	apps "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strings"
)

const	(
	// DescriptionAnnotationKey is annotation key for a description
	DescriptionAnnotationKey = "description"
)

type AppDeploymentSpec struct {
	// Name of the application
	Name string `json:"name"`

	// Docker image path for the application
	ContainerImage string `json:"containerImage"`

	// The name of an image pull secret in case of a private docker repository
	ImagePullSecret *string `json:"imagePullSecret"`

	// Command that is executed instead of container entrypoint, if specified
	ContainerCommand *string `json:"containerCommand"`

	// Arguments for the specified container command or container entrypoint (if command is not specified here)
	ContainerCommandArgs *string `json:"containerCommandArgs"`

	// Number of replicas of the image to maintain
	Replicas int32 `json:"replicas"`

	// Port mappings for the services that is created. The service is created if there is at least
	// one port mapping
	PortMappings []PortMapping `json:"portMappings"`

	// List of user-defined environment variables
	Variables []EnvironmentVariable `json:"variables"`

	// Whether the created service is external
	IsExternal bool `json:"isExternal"`

	// Description of the deployment
	Description *string `json:"description"`

	// Target namespace of the application
	Namespace string `json:"namespace"`

	// Optional memory requirement for the container
	MemoryRequirement *resource.Quantity `json:"memoryRequirement"`

	// Optional CPU requirement for the container
	CpuRequirement *resource.Quantity `json:"cpuRequirement"`

	// Labels that will be defined on Pods/RCs/Services
	Labels []Label `json:"labels"`

	// Whether to run the container as privileged user (essentially equivalent to root on the host)
	RunAsPrivileged bool `json:"runAsPrivileged"`
}

// AppDeploymentFromFileSpec is a specification for deployment from file
type AppDeploymentFromFileSpec struct {
	// Name of the file
	Name string `json:"name"`

	// Namespace that the object should be deployed in
	Namespace string `json:"namespace"`

	// File content
	Content string `json:"content"`

	// Whether validate content before creation or not
	Validate bool `json:"validate"`
}

// AppDeploymentFromFileResponse is a specification for deployment from file
type AppDeploymentFromFileResponse struct {
	// Name of the file
	Name string `json:"name"`

	// File content
	Content string `json:"content"`

	// Error after create resource
	Error string `json:"error"`
}

type PortMapping struct {
	// Port that will be exposed on the service
	Port int32 `json:"port"`

	// Port in the container for the application
	TargetPort int32 `json:"targetPort"`

	// IP protocol for the mapping, e.g., "TCP" or "UDP"
	Protocol v1.Protocol `json:"protocol"`
}

// EnvironmentVariable represents a named variable accessible for containers
type EnvironmentVariable struct {
	// Name of the variable. Must be a C_IDENTIFIER
	Name string `json:"name"`

	// Value of the variable, as defined in kubernetes core API
	Value string `json:"value"`
}

// Label is a structure representing label assignable to Pod/RC/Service
type Label struct {
	// Label key
	Key string `json:"key"`

	// Label value
	Value string `json:"value"`
}

// Protocols is a structure representing supported protocol types for a service
type Protocols struct {
	// Array containing supported protocol types e.g., ["TCP", "UDP"]
	Protocols []v1.Protocol `json:"protocols"`
}

func DeployApp(spec *AppDeploymentSpec, client kubernetes.Interface) error {
	glog.Infof("Deploying %s application into %s namespace", spec.Name, spec.Namespace)

	annotations := map[string]string{}
	if spec.Description != nil {
		annotations[DescriptionAnnotationKey] = *spec.Description
	}

	labels := getLabelsMap(spec.Labels)
	objectMeta := metaV1.ObjectMeta{
		Annotations: annotations,
		Name: spec.Name,
		Labels: labels,
	}

	containerSpec := v1.Container{
		Name: spec.Name,
		Image: spec.ContainerImage,
		SecurityContext: &v1.SecurityContext{
			Privileged: &spec.RunAsPrivileged,
		},
		Resources: v1.ResourceRequirements{
			Requests: make(map[v1.ResourceName]resource.Quantity),
		},
		Env: convertEnvVarsSpec(spec.Variables),
	}

	if spec.ContainerCommand != nil {
		containerSpec.Command = []string{*spec.ContainerCommand}
	}
	if spec.ContainerCommandArgs != nil {
		containerSpec.Args = []string{*spec.ContainerCommandArgs}
	}

	if spec.CpuRequirement != nil {
		containerSpec.Resources.Requests[v1.ResourceCPU] = *spec.CpuRequirement
	}
	if spec.MemoryRequirement != nil {
		containerSpec.Resources.Requests[v1.ResourceMemory] = *spec.MemoryRequirement
	}

	podSpec := v1.PodSpec{
		Containers: []v1.Container{containerSpec},
	}
	if spec.ImagePullSecret != nil {
		podSpec.ImagePullSecrets = []v1.LocalObjectReference{{Name: *spec.ImagePullSecret}}
	}

	podTemplate := v1.PodTemplateSpec{
		ObjectMeta: objectMeta,
		Spec: podSpec,
	}

	deployment := &apps.Deployment{
		ObjectMeta: objectMeta,
		Spec: apps.DeploymentSpec{
			Replicas: &spec.Replicas,
			Template: podTemplate,
			Selector: &metaV1.LabelSelector{
				MatchLabels: labels,
			},
		},
	}
	_, err := client.AppsV1beta2().Deployments(spec.Namespace).Create(deployment)

	if err != nil {
		// TODO(wzt3309) Rollback created resource if it failed
		return err
	}

	if len(spec.PortMappings) > 0 {
		service := &v1.Service{
			ObjectMeta: objectMeta,
			Spec: v1.ServiceSpec{
				Selector: labels,
			},
		}

		if spec.IsExternal {
			service.Spec.Type = v1.ServiceTypeLoadBalancer
		} else {
			service.Spec.Type = v1.ServiceTypeClusterIP
		}

		for _, portMapping := range spec.PortMappings {
			servicePort := v1.ServicePort{
				Name: generatePortMappingName(portMapping),
				Protocol: portMapping.Protocol,
				Port: portMapping.Port,
				TargetPort: intstr.IntOrString{
					Type: intstr.Int,
					IntVal: portMapping.TargetPort,
				},
			}
			service.Spec.Ports = append(service.Spec.Ports, servicePort)
		}

		_, err = client.CoreV1().Services(spec.Namespace).Create(service)
		// TODO(wzt3309) Rollback created resource if it failed
		return err
	}
	return nil
}

// GetAvailableProtocols returns list of available protocols. Currently it is TCP and UDP.
func GetAvailableProtocols() *Protocols {
	return &Protocols{Protocols: []v1.Protocol{v1.ProtocolTCP, v1.ProtocolUDP}}
}

func DeployAppFromFile(cfg *rest.Config, spec *AppDeploymentFromFileSpec) (bool, error) {
	reader := strings.NewReader(spec.Content)
	glog.Infof("Deploy %s file in %s namespace", spec.Name, spec.Namespace)
	d := yaml.NewYAMLOrJSONDecoder(reader, 4096)
	for {
		data := unstructured.Unstructured{}
		if err := d.Decode(&data); err != nil {
			if err == io.EOF {
				return true, nil
			}
			return false, err
		}

		version := data.GetAPIVersion()
		kind := data.GetKind()

		gv, err := schema.ParseGroupVersion(version)
		if err != nil {
			gv = schema.GroupVersion{Version: version}
		}

		groupVersionKind := schema.GroupVersionKind{Group: gv.Group, Version: gv.Version, Kind: kind}

		discoveryClient, err := discovery.NewDiscoveryClientForConfig(cfg)
		if err != nil {
			return false, err
		}

		apiResourceList, err := discoveryClient.ServerResourcesForGroupVersion(version)
		if err != nil {
			return false, err
		}
		apiResources := apiResourceList.APIResources
		var resource *metaV1.APIResource
		for _, apiResource := range apiResources {
			if apiResource.Kind == kind && !strings.Contains(apiResource.Name, "/") {
				resource = &apiResource
				break
			}
		}
		if resource == nil {
			return false, fmt.Errorf("Unknown resource kind: %s", kind)
		}

		dynamicClientPool := dynamic.NewDynamicClientPool(cfg)

		dynamicClient, err := dynamicClientPool.ClientForGroupVersionKind(groupVersionKind)

		if err != nil {
			return false, err
		}

		if strings.Compare(spec.Namespace, "_all") == 0 {
			_, err = dynamicClient.Resource(resource, data.GetNamespace()).Create(&data)
		} else {
			_, err = dynamicClient.Resource(resource, spec.Namespace).Create(&data)
		}

		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func generatePortMappingName(portMapping PortMapping) string {
	return generateName(fmt.Sprintf("%s-%d-%d", strings.ToLower(string(portMapping.Protocol)),
		portMapping.Port, portMapping.TargetPort))
}

func generateName(base string) string {
	maxNameLength := 30
	randomLength := 5
	maxGeneratedNameLength := maxNameLength - randomLength
	if len(base) > maxGeneratedNameLength {
		base = base[:maxGeneratedNameLength]
	}
	return fmt.Sprintf("%s%s", base, rand.String(randomLength))
}

func convertEnvVarsSpec(variables []EnvironmentVariable) []v1.EnvVar {
	var result []v1.EnvVar
	for _, variable := range variables {
		result = append(result, v1.EnvVar{Name: variable.Name, Value: variable.Value})
	}
	return result
}

func getLabelsMap(labels []Label) map[string]string {
	result := make(map[string]string)

	for _, label := range labels {
		result[label.Key] = result[label.Value]
	}

	return result
}