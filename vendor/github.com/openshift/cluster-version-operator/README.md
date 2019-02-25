# Cluster Version Operator (CVO)

## Building and Publishing CVO

```sh
./hack/build-image.sh && REPO=<your personal repo (quay.io/ahbinavdahiya | docker.io/abhinavdahiya)> ./hack/push-image.go
```

1. This builds image locally and then pushes `${VERSION}` and `latest` tags to `${REPO}/origin-cluster-version-operator`.

2. `${VERSION}` encodes the Git commit used to build the images.

## Building release image using local CVO

1. Make sure you have `oc` binary from https://github.com/openshift/origin master as it requires `adm release` subcommand.

2. Run the following command to create release-image at `docker.io/abhinavdahiya/origin-release:latest`:

```sh
oc adm release new -n openshift --server https://api.ci.openshift.org \
    --from-image-stream=origin-v4.0 \
    --to-image-base=docker.io/abhinavdahiya/origin-cluster-version-operator:latest \
    --to-image docker.io/abhinavdahiya/origin-release:latest
```

## Using CVO to render the release-image locally

1. Run the following command to get render the release-payload contents to `/tmp/cvo`

```sh
podman run --rm -ti \
    -v /tmp/cvo:/tmp/cvo:z \
    <release image> \
        render \
        --output-dir=/tmp/cvo \
        --release-image="<release image>"
```

`<release image>` can be personal release image generated using [this](#building-release-image-using-local-cvo) or Origin's release image like `registry.svc.ci.openshift.org/openshift/origin-release:v4.0`.

## Installing CVO and operators in cluster.

1. Use CVO `render` to render all the manifests from release-payload to a directory. [here](#using-cvo-to-render-the-release-payload-locally)

2. Create the operators from the manifests by using `oc create -f <directory when CVO rendered manfiests>`.

## Running CVO tests

```sh
# Run all unit tests
go test ./...

# Run integration tests against a cluster (creates content in a given namespace)
# Requires the CVO CRD to be installed.
export KUBECONFIG=<admin kubeconfig>
TEST_INTEGRATION=1 go test ./... -test.run=^TestIntegration
```