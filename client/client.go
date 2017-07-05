package client

import (
	"github.com/yaronha/kube-crd/crd"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// This file implement all the (CRUD) client methods we need to access our TPR object

func TprClient(cl *rest.RESTClient, namespace string) *tprclient {
	return &tprclient{cl: cl, ns: namespace}
}

type tprclient struct {
	cl *rest.RESTClient
	ns string
}

func (f *tprclient) Create(obj *crd.Example) (*crd.Example, error) {
	var result crd.Example
	err := f.cl.Post().
		Namespace(f.ns).Resource(crd.TPRPlural).
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *tprclient) Update(obj *crd.Example) (*crd.Example, error) {
	var result crd.Example
	err := f.cl.Put().
		Namespace(f.ns).Resource(crd.TPRPlural).
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *tprclient) Delete(name string, options *meta_v1.DeleteOptions) error {
	return f.cl.Delete().
		Namespace(f.ns).Resource(crd.TPRPlural).
		Name(name).Body(options).Do().
		Error()
}

func (f *tprclient) Get(name string) (*crd.Example, error) {
	var result crd.Example
	err := f.cl.Get().
		Namespace(f.ns).Resource(crd.TPRPlural).
		Name(name).Do().Into(&result)
	return &result, err
}

func (f *tprclient) List() (*crd.ExampleList, error) {
	var result crd.ExampleList
	err := f.cl.Get().
		Namespace(f.ns).Resource(crd.TPRPlural).
		Do().Into(&result)
	return &result, err
}

// Create a new List watch for our TPR
func (f *tprclient) NewListWatch() *cache.ListWatch {
	return cache.NewListWatchFromClient(f.cl, crd.TPRPlural, f.ns, fields.Everything())
}
