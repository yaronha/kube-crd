package main

import (
	"fmt"
	"github.com/yaronha/kube-crd/client"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"time"
	"github.com/yaronha/kube-crd/crd"
)

// return rest config, if path not specified assume in cluster config
func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func main() {

	kubeconf := "admin.conf" // Full path to Kube config
	config, err := GetClientConfig(kubeconf)
	if err != nil {
		panic(err.Error())
	}

	// create clientset and create our TPR, this only need to run once
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// note: if the TPR exist our CreateTPR function is set to exit without an error
	err = crd.CreateTPR(clientset)
	if err != nil {
		panic(err)
	}

	// Wait for the TPR to be created before we use it (only needed if its a new one)
	time.Sleep(3 * time.Second)

	// Create a new clientset which include our TPR schema
	tprcs, _, err := crd.NewClient(config)
	if err != nil {
		panic(err)
	}

	// Create a TRP client interface
	tprclient := client.TprClient(tprcs, "default")

	// Create a new Example object and write to k8s
	example := &crd.Example{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:   "example123",
			Labels: map[string]string{"mylabel": "test"},
		},
		Spec: crd.ExampleSpec{
			Foo: "example-text",
			Bar: true,
		},
		Status: crd.ExampleStatus{
			State:   "created",
			Message: "Created, not processed yet",
		},
	}

	_, err = tprclient.Create(example)
	if err != nil {
		panic(err)
	}

	// List all Example objects
	items, err := tprclient.List()
	if err != nil {
		panic(err)
	}
	fmt.Printf("List:\n%s\n", items)

	// Watch for changes in Example objects and fire Add, Delete, Update callbacks (i.e. Controller)
	_, controller := cache.NewInformer(
		tprclient.NewListWatch(),
		&crd.Example{},
		time.Minute*10,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				fmt.Printf("add: %s \n", obj)
			},
			DeleteFunc: func(obj interface{}) {
				fmt.Printf("delete: %s \n", obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				fmt.Printf("Update old: %s \n      New: %s\n", oldObj, newObj)
			},
		},
	)

	stop := make(chan struct{})
	go controller.Run(stop)

	// Wait forever
	select {}
}
