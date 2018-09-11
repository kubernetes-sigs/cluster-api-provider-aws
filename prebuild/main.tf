provider "aws" {
  region = "${var.aws_region}"
}

data "aws_availability_zones" "azs" {}
data "aws_region" "current" {}

module "vpc" {
  source          = "terraform-aws-modules/vpc/aws"
  name            = "${var.vpc_name}"
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
    Owner       = "user"
    Environment = "dev"
  }

  vpc_tags = {
    Name = "${var.vpc_name}"
  }
}

# Generate Manifest Dir
resource "template_dir" "manifests" {
  source_dir      = "${path.module}/resources"
  destination_dir = "${path.cwd}/tfManifests"

  vars {
    aws_availability_zone        = "${data.aws_availability_zones.azs.names[0]}"
    aws_machine_controller_image = "${var.container_images["aws_machine_controller"]}"
    aws_region                   = "${data.aws_region.current.name}"
    cluster_apiserver_image      = "${var.container_images["cluster_apiserver"]}"
    cluster_cidr                 = "${var.vpc_cidr}"
    cluster_name                 = "${var.cluster_name}"
    cluster_domain               = "${var.cluster_domain}"
    cluster_namespace            = "${var.cluster_namespace}"
    cluster_security_group       = "${aws_security_group.cluster_default.id}"
    controller_manager_image     = "${var.container_images["controller_manager"]}"
    etcd_image                   = "${var.container_images["etcd"]}"
    apiserver_image              = "${var.container_images["cluster_apiserver"]}"
    pod_cidr                     = "${var.vpc_private_networks[1]}"
    service_cidr                 = "${var.vpc_private_networks[0]}"
    ssh_key_name                 = "${var.sshKey}"
  }
}
