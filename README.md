# Kubernetes Custom Resources (CRD) Tutorial

Tutorial for building Kubernetes Custom Resources (CRD) extensions
you can see the full tutorial documentation in: [The New Stack](https://thenewstack.io/extend-kubernetes-1-7-custom-resources)

**Note:** CustomResourceDefinition (CRD) is the successor of the deprecated ThirdPartyResource.

this example is based on Kubernetes [apiextensions-apiserver](https://github.com/kubernetes/apiextensions-apiserver) example  

## Organization 

the example contain 3 files:

* crd      - define and register our CRD class 
* client   - client library to create and use our CRD (CRUD)
* kube-crd - main part, demonstrate how to create, use, and watch our CRD

## Running

```
# assumes you have a working kubeconfig, not required if operating in-cluster
go run *.go -kubeconf=$HOME/.kube/config
```


## kube-crd

kube-crd demonstrates the CRD usage, it shoes how to:

1. Connect to the Kubernetes cluster 
2. Create the new CRD if it doesn't exist  
3. Create a new custom client 
4. Create a new Example object using the client library we created 
5. Create a controller that listens to events associated with new resources

The example CRD is in the following structure:


```go
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
```

* The Metadata part contain standard Kubernetes properties like name, namespace, labels, and annotations 
* The Spec contain the desired resource configuration 
* The Status part is usually filled by the controller in response to Spec updates 

