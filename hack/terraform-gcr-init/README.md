terraform-gcr-init
==================

This can be used as a Terraform module to set up a Google Project
for each developer working on cluster-api-provider-aws as a place for them
to push container images.

**NOTE:**
> To make GCR usable, it will make the repositories publicly readable.

## Example use

``` hcl
module "a_dev" {
  source  = "<path to here>"

  // The name of the developer
  short_name = "a_dev"

  // A Google entity that has ownership of the project
  owners = "group:cluster-provider-aws@sigs.k8s.io"

  // The organisation ID
  org_id = "123456789012"

  // The billing account for this project
  billing_account = "012345-6789AB-CDEF01"
}
```
