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
package client

import (
	"github.com/yaronha/kube-crd/crd"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// This file implement all the (CRUD) client methods we need to access our TPR object

func CrdClient(cl *rest.RESTClient, namespace string) *crdclient {
	return &crdclient{cl: cl, ns: namespace, plural: crd.CRDPlural}
}

type crdclient struct {
	cl     *rest.RESTClient
	ns     string
	plural string
}

func (f *crdclient) Create(obj *crd.Example) (*crd.Example, error) {
	var result crd.Example
	err := f.cl.Post().
		Namespace(f.ns).Resource(f.plural).
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *crdclient) Update(obj *crd.Example) (*crd.Example, error) {
	var result crd.Example
	err := f.cl.Put().
		Namespace(f.ns).Resource(f.plural).
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *crdclient) Delete(name string, options *meta_v1.DeleteOptions) error {
	return f.cl.Delete().
		Namespace(f.ns).Resource(f.plural).
		Name(name).Body(options).Do().
		Error()
}

func (f *crdclient) Get(name string) (*crd.Example, error) {
	var result crd.Example
	err := f.cl.Get().
		Namespace(f.ns).Resource(f.plural).
		Name(name).Do().Into(&result)
	return &result, err
}

func (f *crdclient) List() (*crd.ExampleList, error) {
	var result crd.ExampleList
	err := f.cl.Get().
		Namespace(f.ns).Resource(f.plural).
		Do().Into(&result)
	return &result, err
}

// Create a new List watch for our TPR
func (f *crdclient) NewListWatch() *cache.ListWatch {
	return cache.NewListWatchFromClient(f.cl, f.plural, f.ns, fields.Everything())
}
