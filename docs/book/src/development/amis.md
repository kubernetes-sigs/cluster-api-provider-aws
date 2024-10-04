# Publish AMIs

Publishing new AMIs is currently a manual process but it will be automated in th every near future (see [this issue](https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1982) for progress).

## Pre-reqs

- You must have admin access to the CNCF AWAS account used for the AMIs (819546954734)

## Process

1. Clone [image-builder](https://github.com/kubernetes-sigs/image-builder)
2. Open a terminal
3. Set the AWS environment variables for the AMI account
4. Change directory into `images/capi`
5. Install dependencies by running:

```shell
make deps-ami
```

6. Build the AMIs using:

```shell
make build-ami-ubuntu-2004
make build-ami-ubuntu-2204
make build-ami-ubuntu-2404
make build-ami-flatcar
make build-ami-rhel-8
```
> NOTE: there are some issues with the RHEL based images at present.
