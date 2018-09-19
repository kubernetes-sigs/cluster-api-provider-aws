// Your aws account user name
variable "aws_user" {
  type = "string"
}

variable "aws_region" {
  type    = "string"
  default = "us-east-1"
}

variable "vpc_cidr" {
  type    = "string"
  default = "10.0.0.0/16"
}

variable "vpc_public_networks" {
  default = [
    "10.0.101.0/24",
    "10.0.102.0/24",
    "10.0.103.0/24",
  ]
}

variable "vpc_private_networks" {
  default = [
    "10.0.1.0/24",
    "10.0.2.0/24",
    "10.0.3.0/24",
  ]
}

variable "environment_id" {
  type        = "string"
  default     = "testCluster"
  description = "The id of the environment."
}
