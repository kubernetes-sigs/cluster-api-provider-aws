provider "aws" {
  region = "${var.aws_region}"
}

data "aws_availability_zones" "azs" {}
data "aws_region" "current" {}

module "vpc" {
  source          = "terraform-aws-modules/vpc/aws"
  name            = "vpc-${var.environment_id}"
  cidr            = "${var.vpc_cidr}"
  azs             = "${data.aws_availability_zones.azs.names}"
  public_subnets  = "${var.vpc_public_networks}"
  private_subnets = "${var.vpc_private_networks}"

  public_subnet_tags = {
    Name = "${var.environment_id}"
  }

  private_subnet_tags = {
    Name = "${var.environment_id}"
  }

  enable_dns_hostnames = true
  enable_nat_gateway = false
  single_nat_gateway = false

  public_subnet_tags = {
    Name = "${var.environment_id}-worker-foo"
  }

  tags = {
    Owner       = "${var.aws_user}"
    Environment = "dev"
  }

  vpc_tags = {
    Name = "vpc-${var.environment_id}"
  }
}

resource "aws_iam_instance_profile" "test_profile" {
  name = "${var.environment_id}-worker-profile"
  role = "${aws_iam_role.role.name}"
}

resource "aws_iam_role" "role" {
  name = "${var.environment_id}-role"
  path = "/"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
               "Service": "ec2.amazonaws.com"
            },
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}
EOF
}

output "vpc_id" {
  value       = "${module.vpc.vpc_id}"
  description = "The ID of the VPC"
}
