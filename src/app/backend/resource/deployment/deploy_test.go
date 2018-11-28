package deployment

import (
	apps "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"
	"reflect"
	"testing"
)

func TestDeployApp(t *testing.T) {
	replicas := int32(0)
	namespace := "foo-namespace"
	spec := &AppDeploymentSpec{
		Namespace:       namespace,
		Name:            "foo-name",
		RunAsPrivileged: true,
	}

	expected := &apps.Deployment{
		ObjectMeta: metaV1.ObjectMeta{
			Name: "foo-name",
			Labels: map[string]string{},
			Annotations: map[string]string{},
		},
		Spec: apps.DeploymentSpec{
			Selector: &metaV1.LabelSelector{
				MatchLabels: map[string]string{},
			},
			Replicas: &replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Name: "foo-name",
					Labels: map[string]string{},
					Annotations: map[string]string{},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "foo-name",
							SecurityContext: &v1.SecurityContext{
								Privileged: &spec.RunAsPrivileged,
							},
							Resources: v1.ResourceRequirements{
								Requests: make(map[v1.ResourceName]resource.Quantity),
							},
						},
					},
				},
			},
		},
	}

	testClient := fake.NewSimpleClientset()

	DeployApp(spec, testClient)

	createAction := testClient.Actions()[0].(core.CreateActionImpl)
	if len(testClient.Actions()) != 1 {
		t.Errorf("Expected one create action but got %#v", len(testClient.Actions()))
	}

	if createAction.GetNamespace() != namespace {
		t.Errorf("Expected namespace to be %#v but got %#v", namespace, createAction.GetNamespace())
	}

	deployment := createAction.GetObject().(*apps.Deployment)
	if !reflect.DeepEqual(deployment, expected) {
		t.Errorf("Expected deployment \n%#v\n to be created but got \n%#v\n",
			expected, deployment)
	}
}
