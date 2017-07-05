package crd

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	v1b1e "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/rest"
)

const (
	TPRName        string = "example"
	TPRPlural      string = "examples"
	TPRGroup       string = "myorg.io"
	TPRVersion     string = "v1"
	TPRDescription string = "My Example"
)

// Create the TPR resource, ignore error if it already exists
func CreateTPR(clientset kubernetes.Interface) error {
	tpr := &v1b1e.ThirdPartyResource{
		ObjectMeta:  meta_v1.ObjectMeta{Name: TPRName + "." + TPRGroup},
		Versions:    []v1b1e.APIVersion{{Name: TPRVersion}},
		Description: TPRDescription,
	}

	_, err := clientset.Extensions().ThirdPartyResources().Create(tpr)
	if err != nil && apierrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}

// Definition of our TPR Example class
type Example struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec               ExampleSpec   `json:"spec"`
	Status             ExampleStatus `json:"status,omitempty"`
}
type ExampleSpec struct {
	Foo string `json:"foo"`
	Bar bool   `json:"bar"`
	Baz int    `json:"baz,omitempty"`
}

type ExampleStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}

type ExampleList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []Example `json:"items"`
}

// Create a  Rest client with the new TPR Schema
var SchemeGroupVersion = schema.GroupVersion{Group: TPRGroup, Version: TPRVersion}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Example{},
		&ExampleList{},
	)
	meta_v1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func NewClient(cfg *rest.Config) (*rest.RESTClient, *runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	SchemeBuilder := runtime.NewSchemeBuilder(addKnownTypes)
	if err := SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}
	config := *cfg
	config.GroupVersion = &SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{
		CodecFactory: serializer.NewCodecFactory(scheme)}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, nil, err
	}
	return client, scheme, nil
}
