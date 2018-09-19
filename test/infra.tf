provider "aws" {
  region = "us-east-1"
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "${var.enviroment_id}"

  cidr = "10.0.0.0/16"

  azs = [
    "us-east-1a",
    "us-east-1b",
    "us-east-1c",
  ]

  private_subnets = [
    "10.0.1.0/24",
    "10.0.2.0/24",
    "10.0.3.0/24",
  ]

  public_subnets = [
    "10.0.101.0/24",
    "10.0.102.0/24",
    "10.0.103.0/24",
  ]

  enable_nat_gateway = true
  single_nat_gateway = true
  enable_dns_hostnames = true

  public_subnet_tags = {
    Name = "${var.enviroment_id}-worker-foo"
  }

  tags = {
    Owner       = "user"
    Environment = "dev"
  }

  vpc_tags = {
    Name = "${var.enviroment_id}"
  }
}

// Security group to verify the machineSet filter by tag works
resource "aws_security_group" "test" {
  name   = "${var.enviroment_id}-sg"
  vpc_id = "${module.vpc.vpc_id}"

  tags = {
    Name = "${var.enviroment_id}_worker_sg"
  }
}

resource "aws_security_group_rule" "test" {
  type              = "ingress"
  security_group_id = "${aws_security_group.test.id}"

  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 0
  to_port     = 0
}

resource "aws_security_group_rule" "allow_all" {
  type              = "ingress"
  security_group_id = "${aws_security_group.test.id}"

  cidr_blocks = ["0.0.0.0/0"]
  from_port       = 22
  to_port         = 22
  description     = "SSH"
  protocol        = "tcp"
}

// Instance profile to verify the machineSet iamInstanceProfile field works
resource "aws_iam_instance_profile" "test_profile" {
  name = "${var.enviroment_id}-worker-profile"
  role = "${aws_iam_role.role.name}"
}

resource "aws_iam_role" "role" {
  name = "${var.enviroment_id}-role"
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
