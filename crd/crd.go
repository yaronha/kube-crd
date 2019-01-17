/*
Copyright 2016 Iguazio Systems Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package crd

import (
	"reflect"

	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextcs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

const (
	CRDPlural      string = "examples"
	CRDGroup       string = "myorg.io"
	CRDVersion     string = "v1"
	FullCRDName    string = CRDPlural + "." + CRDGroup
)

// Create the CRD resource, ignore error if it already exists
func CreateCRD(clientset apiextcs.Interface) error {
	crd := &apiextv1beta1.CustomResourceDefinition{
		ObjectMeta: meta_v1.ObjectMeta{Name: FullCRDName},
		Spec: apiextv1beta1.CustomResourceDefinitionSpec{
			Group:   CRDGroup,
			Version: CRDVersion,
			Scope:   apiextv1beta1.NamespaceScoped,
			Names:   apiextv1beta1.CustomResourceDefinitionNames{
				Plural: CRDPlural,
				Kind:   reflect.TypeOf(Example{}).Name(),
			},
		},
	}

	_, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil && apierrors.IsAlreadyExists(err) {
		return nil
	}
	return err

	// Note the original apiextensions example adds logic to wait for creation and exception handling
}

// Definition of our CRD Example class
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

func (in *Example) DeepCopyInto(out *Example) {
    out.TypeMeta = in.TypeMeta
    out.ObjectMeta = in.ObjectMeta
    out.Spec = ExampleSpec{
		Foo: in.Spec.Foo,
		Bar: in.Spec.Bar,
		Baz: in.Spec.Baz,
	}
	out.Status = ExampleStatus{
		State: in.Status.State,
		Message: in.Status.Message,
    }
}

// DeepCopyObject returns a generically typed copy of an object
func (in *Example) DeepCopyObject() runtime.Object {
    out := Example{}
    in.DeepCopyInto(&out)

    return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *ExampleList) DeepCopyObject() runtime.Object {
    out := ExampleList{}
    out.TypeMeta = in.TypeMeta
    out.ListMeta = in.ListMeta

    if in.Items != nil {
        out.Items = make([]Example, len(in.Items))
        for i := range in.Items {
            in.Items[i].DeepCopyInto(&out.Items[i])
        }
    }

    return &out
}

// Create a  Rest client with the new CRD Schema
var SchemeGroupVersion = schema.GroupVersion{Group: CRDGroup, Version: CRDVersion}

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
