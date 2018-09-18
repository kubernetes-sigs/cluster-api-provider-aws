provider "aws" {
  region = "${var.aws_region}"
}

data "aws_availability_zones" "azs" {}
data "aws_region" "current" {}

module "vpc" {
  source          = "terraform-aws-modules/vpc/aws"
  name            = "vpc-${var.cluster_name}"
  cidr            = "${var.vpc_cidr}"
  azs             = "${data.aws_availability_zones.azs.names}"
  public_subnets  = "${var.vpc_public_networks}"
  private_subnets = "${var.vpc_private_networks}"

  public_subnet_tags = {
    Name = "${var.cluster_name}"
  }

  private_subnet_tags = {
    Name = "${var.cluster_name}"
  }

  enable_nat_gateway = true
  single_nat_gateway = true

  tags = {
    Owner       = "${var.aws_user}"
    Environment = "dev"
  }

  vpc_tags = {
    Name = "vpc-${var.cluster_name}"
  }
}
