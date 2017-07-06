# Kubernetes TPR Tutorial

Tuturial for building a Kubernetes Custom Resources (CRD) extension
you can see the full tutorial in: TBD

# Organization 

the example contain 3 files:

1. crd      - define and register our TPR class 
2. client   - client library to create and use our CRD (CRUD)
3. kube-crd - main part, demonstrate how to create, use, and watch our CRD

# kube-crd

kube-crd demonstrates the CRD usage, it shoes how to:

1. Connect to the Kubernetes cluster 
2. Create the new CRD if it doesn't exist  
3. Create a new custom client 
4. Create a new Example object using the client library we created 
5. Create a controller that listens to events associated with new resources

