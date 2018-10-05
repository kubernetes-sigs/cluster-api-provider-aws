// Example: "adev"
variable "short_name" {
  type = "string"
}

// Example: "group:cluster-provider-aws@sigs.k8s.io"
variable "owners" {
  type = "string"
}

// Example: 123456789012
variable "org_id" {
  type = "string"
}

// Example: 012345-6789AB-CDEF01"
variable "billing_account" {
  type = "string"
}
