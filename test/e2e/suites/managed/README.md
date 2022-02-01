# Managed (i.e. EKS) E2E Test Notes

Due to the time in takes (15-20 mins) to provision an EKS cluster in AWS the e2e tests for the managed side of CAPA have been implemented in a way to reduce the number of clusters created.

Where possible we re-use an EKS cluster and have tried not to create clusters in a `BeforeEach`. Instead we create a cluster in steps and at each step perform tests.

For example, in [eks_test.go](eks_test.go) we perform the following steps:

1. Apply an AWSManagedControlPlane to create a EKS cluster without any nodes
2. Perform tests against the control plane
3. Apply a MachineDeployment
4. Perform tests against the machine deployment
5. Apply a AWSManagedMachinePool
6. Perform tests against the machine pool