# Using Packer and baking AMIs

## Overview

[Packer](http://packer.io/) is well known tool for baking images of any kind.
We use it to bake our AMIs.

## Prerequisites

* `packer` binary
* [packer-goss](https://github.com/YaleUniversity/packer-provisioner-goss) plugin
* ansible

## Plugin instalation

To install `packer-goss` plugin the following should be executed inside of the
`build/amis/packer` directory:

```bash
$ curl -o packer-goss  https://github.com/YaleUniversity/packer-provisioner-goss/releases/download/v0.3.0/packer-provisioner-goss-v0.3.0-linux-amd64

$ chmod +x packer-goss
```

## Running Packer

The following command should build all the AMIs:

```bash
$ AWS_REGION=us-east-1 packer build -var-file=base-images-us-east-1.json packer.json
```

**NOTE** that AWS credentials have to be set.