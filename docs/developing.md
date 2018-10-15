# Dependencies

Instead of checking vendor code into this repo we are only checking the files required to rebuild the vendor directory. Modules to be determined.

To get started run `dep ensure`.

## clusterctl

[clusterctl](https://github.com/kubernetes-sigs/cluster-api/tree/master/clusterctl#getting-started) is a bit of a strange beast. You must compile this from the vendor directory in this repo.

## Building

### Kubebuilder

Kubebuilder should be available and all binaries associated should be in your path. Please see installation instructions https://book.kubebuilder.io/quick_start.html

If you update any of the types you must regenerate the CRDs with `make crds`.

