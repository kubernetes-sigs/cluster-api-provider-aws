manifests-gen
=============

This tool is used to generate manifests for Cluster API Providers suitable for use within OpenShift Container Platform.

The tool makes the following assumptions, and uses relative paths based on those:

    * The forked provider repository has directories `openshift` and `openshift/tools`
    * `manifests-gen`'s working directory is `openshift/tools` (invoked via `openshift/Makefile`)

Most often, the base Cluster API Provider YAML manifests will be stored in `config/default`.

If there are any OpenShift-specific changes that do not warrant inclusion into this tool,
the directory can be overriden with `--kustomize-dir`. Often, this should be `--kustomize-dir=openshift`
and a `openshift/kustomization.yaml` file should be present that references `config/default` as the base.

Changes that warrant making edits via this tool are things that apply to all or many resources.
Changes that target a single resource, such as changing a single, specific `Deployment`'s
container arguments should be applies as Kustomize patches local to the provider repo.
