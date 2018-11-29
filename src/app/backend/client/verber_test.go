package client

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestResourceVerber_Put(t *testing.T) {
	testClient := fake.NewSimpleClientset()
	testVerber := NewResourceVerber(testClient.CoreV1().RESTClient(),
		testClient.ExtensionsV1beta1().RESTClient(), testClient.AppsV1beta2().RESTClient(),
		testClient.BatchV1().RESTClient(), testClient.BatchV1beta1().RESTClient(), testClient.StorageV1().RESTClient())

	//reqObj := &v1.Pod{ObjectMeta: metaV1.ObjectMeta{Name: "foo-name"}}
	//reqBodyExpected, err := runtime.Encode(scheme.Codecs.LegacyCodec(v1.SchemeGroupVersion), reqObj)
	//if err != nil {
	//	t.Errorf("unexpected error: %v", err)
	//}

	object := &runtime.Unknown{
		TypeMeta: runtime.TypeMeta{
			APIVersion: "v1",
			Kind: "Pod",
		},
		//Raw: reqBodyExpected,
	}
	testVerber.Put("pods", true,"default", "foo-name", object)
}
