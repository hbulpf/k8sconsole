package container

import (
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PodContainerList is a list of containers of a pod
type PodContainerList struct {
	Containers []string `json:"containers"`
}

// GetPodContainers returns containers that a
func GetPodContainers(client kubernetes.Interface, namespace, name string) (*PodContainerList, error) {
	pod, err := client.CoreV1().Pods(namespace).Get(name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	containers := &PodContainerList{Containers: make([]string, 0)}

	for _, container := range pod.Spec.Containers {
		containers.Containers = append(containers.Containers, container.Name)
	}

	return containers, nil
}

