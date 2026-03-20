# Managed (i.e. EKS) E2E Test Notes

Due to the time it takes (15-20 mins) to provision an EKS cluster in AWS the e2e tests for the managed side of CAPA have been implemented in a way to reduce the number of clusters created.

Where possible we re-use an EKS cluster and have tried not to create clusters in a `BeforeEach`. Instead we create a cluster in steps and at each step perform tests.

For example, in [eks_test.go](eks_test.go) we perform the following steps:

1. Apply an AWSManagedControlPlane to create a EKS cluster without any nodes
2. Perform tests against the control plane
3. Apply a MachineDeployment
4. Perform tests against the machine deployment
5. Apply a AWSManagedMachinePool
6. Perform tests against the machine pool
7. Apply a AWSMachinePool
8. Perform tests against the machine pool

## EKS Quick Test

The [eks_quick_test.go](eks_quick_test.go) implements a minimalist quick test for EKS clusters that is designed to run on every PR as a sanity check. This test is tagged with `[PR-Blocking]` and `[smoke]` labels.

This test creates a basic EKS cluster with:
1. An EKS control plane
2. A single managed node group

The test validates basic cluster provisioning and then performs cleanup. It is designed to be fast and provide quick feedback on PRs without running the full test suite.