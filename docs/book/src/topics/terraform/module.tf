variable "cluster_name" {
  type = string
  default = "new-test-aws"
}

variable "vpc_cidr" {
  type = string
  default = "10.0.0.0/16"
}

variable "zone_1" {
  type = string
  default = "eu-west-1a"
}

variable "zone_2" {
  type = string
  default = "eu-west-1b"
}

variable "zone_3" {
  type = string
  default = "eu-west-1b"
}

locals {
  cluster_tag = join("",["kubernetes.io/cluster/", var.cluster_name])
  role_tag = "sigs.k8s.io/cluster-api-provider-aws/role"
  public_elb_role_tag = "kubernetes.io/role/elb"
  internal_elb_role_tag = "kubernetes.io/role/internal-elb"
  public_subnet_prefix = join("-",[var.cluster_name, "subnet-public"])
  private_subnet_prefix = join("-",[var.cluster_name, "subnet-private"])
  public_rt_prefix = join("-", [var.cluster_name, "rt", "public"])
  private_rt_prefix = join("-", [var.cluster_name, "rt", "private"])
  sg_description_prefix = join("",["Kubernetes cluster ", var.cluster_name, ": "])
  rt_prefix = join("-", [var.cluster_name, "rt"])
}

resource "aws_vpc" "cluster" {
  cidr_block = var.vpc_cidr
  tags = {
    (local.cluster_tag) = "shared"
    Name = join("-",[var.cluster_name, "vpc"])
    (local.role_tag) = "common"
  }
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.cluster.id

  tags = {
    (local.cluster_tag) = "shared"
    (local.role_tag) = "common"
    Name = join("-",[var.cluster_name, "igw"])
  }
}

resource "aws_subnet" "public-1" {
  // test-aws-internal-subnet-public-eu-west-1a
  cidr_block = "10.0.0.0/20"
  availability_zone = var.zone_1
  vpc_id = aws_vpc.cluster.id
  map_public_ip_on_launch = true
  tags = {
    (local.cluster_tag) = "shared"
    Name = join("-",[local.public_subnet_prefix, var.zone_1])
    (local.public_elb_role_tag) = "1"
    (local.role_tag) = "public"
  }
}

resource "aws_subnet" "public-2" {
  // test-aws-internal-subnet-public-eu-west-1a
  cidr_block = "10.0.32.0/20"
  availability_zone = var.zone_2
  vpc_id = aws_vpc.cluster.id
  map_public_ip_on_launch = true
  tags = {
    (local.cluster_tag) = "shared"
    Name = join("-",[local.public_subnet_prefix, var.zone_2])
    (local.public_elb_role_tag) = "1"
    (local.role_tag) = "public"
  }
}

resource "aws_subnet" "public-3" {
  cidr_block = "10.0.16.0/20"
  availability_zone = var.zone_3
  vpc_id = aws_vpc.cluster.id
  map_public_ip_on_launch = true
  tags = {
    (local.cluster_tag) = "shared"
    Name = join("-",[local.public_subnet_prefix, var.zone_3])
    (local.public_elb_role_tag) = "1"
    (local.role_tag) = "public"
  }
}

resource "aws_eip" "natgw-1" {
  vpc      = true
  tags = {
    (local.cluster_tag) = "shared"
    Name = join("-",[var.cluster_name, "eip", "apiserver"])
    (local.role_tag) = "apiserver"
  }
}

resource "aws_eip" "natgw-2" {
  vpc      = true
  tags = {
    (local.cluster_tag) = "shared"
    Name = join("-",[var.cluster_name, "eip", "apiserver"])
    (local.role_tag) = "apiserver"
  }
}

resource "aws_eip" "natgw-3" {
  vpc      = true
  tags = {
    (local.cluster_tag) = "shared"
    Name = join("-",[var.cluster_name, "eip", "apiserver"])
    (local.role_tag) = "apiserver"
  }
}

resource "aws_nat_gateway" "natgw-1" {
  allocation_id = aws_eip.natgw-1.id
  subnet_id     = aws_subnet.public-1.id
  tags = {
    (local.cluster_tag) = "shared"
    Name = join("-",[var.cluster_name, "nat"])
    (local.public_elb_role_tag) = "1"
    (local.role_tag) = "common"
  }
}

resource "aws_nat_gateway" "natgw-2" {
  allocation_id = aws_eip.natgw-2.id
  subnet_id     = aws_subnet.public-2.id
  tags = {
    (local.cluster_tag) = "shared"
    Name = join("-",[var.cluster_name, "nat"])
    (local.public_elb_role_tag) = "1"
    (local.role_tag) = "common"
  }
}

resource "aws_nat_gateway" "natgw-3" {
  allocation_id = aws_eip.natgw-3.id
  subnet_id     = aws_subnet.public-3.id
  tags = {
    (local.cluster_tag) = "shared"
    Name = join("-",[var.cluster_name, "nat"])
    (local.public_elb_role_tag) = "1"
    (local.role_tag) = "common"
  }
}

resource "aws_route_table" "public-1" {
  vpc_id = aws_vpc.cluster.id

  route {
    cidr_block        = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    Name = join("-", [local.public_rt_prefix, var.zone_1])
    (local.cluster_tag) = "shared"
    (local.role_tag) = "common"
  }
}

resource "aws_route_table_association" "public-1" {
  subnet_id      = aws_subnet.public-1.id
  route_table_id = aws_route_table.public-1.id
}

resource "aws_route_table" "public-2" {
  vpc_id = aws_vpc.cluster.id

  route {
    cidr_block        = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    Name = join("-", [local.public_rt_prefix, var.zone_2])
    (local.cluster_tag) = "shared"
    (local.role_tag) = "common"
  }
}

resource "aws_route_table_association" "public-2" {
  subnet_id      = aws_subnet.public-2.id
  route_table_id = aws_route_table.public-2.id
}

resource "aws_route_table" "public-3" {
  vpc_id = aws_vpc.cluster.id

  route {
    cidr_block        = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    Name = join("-", [local.public_rt_prefix, var.zone_3])
    (local.cluster_tag) = "shared"
    (local.role_tag) = "common"
  }
}

resource "aws_route_table_association" "public-3" {
  subnet_id      = aws_subnet.public-3.id
  route_table_id = aws_route_table.public-3.id
}

resource "aws_route_table" "private-1" {
  vpc_id = aws_vpc.cluster.id

  route {
    cidr_block        = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.natgw-1.id
  }

  tags = {
    Name = join("-", [local.private_rt_prefix, var.zone_1])
    (local.cluster_tag) = "shared"
    (local.role_tag) = "common"
  }
}

resource "aws_route_table_association" "private-1" {
  subnet_id      = aws_subnet.private-1.id
  route_table_id = aws_route_table.private-1.id
}

resource "aws_route_table" "private-2" {
  vpc_id = aws_vpc.cluster.id

  route {
    cidr_block        = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.natgw-2.id
  }

  tags = {
    Name = join("-", [local.private_rt_prefix, var.zone_2])
    (local.cluster_tag) = "shared"
    (local.role_tag) = "common"
  }
}

resource "aws_route_table_association" "private-2" {
  subnet_id      = aws_subnet.private-2.id
  route_table_id = aws_route_table.private-2.id
}

resource "aws_route_table" "private-3" {
  vpc_id = aws_vpc.cluster.id

  route {
    cidr_block        = "0.0.0.0/0"
    nat_gateway_id = aws_nat_gateway.natgw-3.id
  }

  tags = {
    Name = join("-", [local.private_rt_prefix, var.zone_3])
    (local.cluster_tag) = "shared"
    (local.role_tag) = "common"
  }
}

resource "aws_route_table_association" "private-3" {
  subnet_id      = aws_subnet.private-3.id
  route_table_id = aws_route_table.private-3.id
}

resource "aws_subnet" "private-1" {
  cidr_block = "10.0.64.0/20"
  availability_zone = var.zone_1
  vpc_id = aws_vpc.cluster.id
  tags = {
    (local.cluster_tag) = "shared"
    Name = join("-",[local.private_subnet_prefix, var.zone_1])
    (local.internal_elb_role_tag) = "1"
    (local.role_tag) = "private"
  }
}


resource "aws_subnet" "private-2" {
  cidr_block = "10.0.128.0/20"
  availability_zone = var.zone_2
  vpc_id = aws_vpc.cluster.id
  tags = {
    (local.cluster_tag) = "shared"
    Name = join("-",[local.private_subnet_prefix, var.zone_2])
    (local.internal_elb_role_tag) = "1"
    (local.role_tag) = "private"
  }
}

resource "aws_subnet" "private-3" {
  cidr_block = "10.0.192.0/20"
  availability_zone = var.zone_3
  vpc_id = aws_vpc.cluster.id
  tags = {
    (local.cluster_tag) = "shared"
    Name = join("-",[local.private_subnet_prefix, var.zone_3])
    (local.internal_elb_role_tag) = "1"
    (local.role_tag) = "private"
  }
}

resource "aws_security_group" "bastion" {
  name        = join("-",[var.cluster_name,"bastion"])
  description = join("",[local.sg_description_prefix, "bastion"])
  vpc_id      = (aws_vpc.cluster.id)


  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = join("",[var.cluster_name,"-bastion"])
    (local.cluster_tag) = "shared"
    (local.role_tag) = "bastion"
  }
}

resource "aws_security_group" "node" {
  name        = join("-",[var.cluster_name,"node"])
  description = join("",[local.sg_description_prefix, "node"])
  vpc_id      = (aws_vpc.cluster.id)

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = join("",[var.cluster_name,"-node"])
    (local.cluster_tag) = "shared"
    (local.role_tag) = "node"
  }
}

resource "aws_security_group_rule" "node_port_services" {
  description = "Node Port Services"
  type              = "ingress"
    from_port   = 30000
    to_port     = 32767
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    security_group_id = (aws_security_group.node.id)
}

resource "aws_security_group_rule" "bgp" {
  type              = "ingress"
  description ="bgp (calico) node"
    from_port   = 179
    to_port     = 179
    protocol    = "tcp"
    self = true
    security_group_id = (aws_security_group.node.id)
}

resource "aws_security_group_rule" "ipip" {
  type              = "ingress"
  description ="IP-in-IP (calico) node"
    protocol    = 4
    from_port = 0
    to_port = 0
    self = true
    security_group_id = (aws_security_group.node.id)
}

resource "aws_security_group_rule" "ssh" {
  description = "SSH"
  type              = "ingress"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    source_security_group_id = (aws_security_group.bastion.id)
    security_group_id = (aws_security_group.node.id)
}

resource "aws_security_group_rule" "kubelet" {
  type              = "ingress"
    description = "Kubelet API"
    from_port   = 10250
    to_port     = 10250
    protocol    = "tcp"
    self = true
  security_group_id = (aws_security_group.node.id)
}

resource "aws_security_group_rule" "controlplane-kubelet-api" {
  type              = "ingress"
  to_port           = 10250
  protocol          = "tcp"
  from_port         = 10250
  security_group_id = (aws_security_group.node.id)
  source_security_group_id = (aws_security_group.controlplane.id)
}

resource "aws_security_group_rule" "controlplane-bgp" {
  description ="bgp (calico)"
  type              = "ingress"
  to_port           = 179
  protocol          = "tcp"
  from_port         = 179
  security_group_id = (aws_security_group.node.id)
  source_security_group_id = (aws_security_group.controlplane.id)
}

resource "aws_security_group_rule" "controlplane-ipip" {
  type              = "ingress"
description ="IP-in-IP (calico)"
  protocol          = 4
  from_port = 0
  to_port = 0
  security_group_id = (aws_security_group.node.id)
  source_security_group_id = (aws_security_group.controlplane.id)
}

resource "aws_security_group" "apiserver-lb" {
  name        = join("-",[var.cluster_name,"apiserver-lb"])
  description = join("",[local.sg_description_prefix, "apiserver-lb"])
  vpc_id      = (aws_vpc.cluster.id)


  ingress {
    description = "Kubernetes API"
    from_port   = 6443
    to_port     = 6443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = join("-",[var.cluster_name,"apiserver-lb"])
    (local.cluster_tag) = "shared"
    (local.role_tag) = "apiserver-lb"
  }
}

resource "aws_security_group" "controlplane" {
  name        = join("-",[var.cluster_name,"controlplane"])
  description = join("",[local.sg_description_prefix, "controlplane"])
  vpc_id      = aws_vpc.cluster.id

  ingress {
    description = "Kubernetes API"
    from_port   = 6443
    to_port     = 6443
    protocol    = "tcp"
    security_groups = [aws_security_group.node.id]
  }

  ingress {
    description = "Kubernetes API"
    from_port   = 6443
    to_port     = 6443
    protocol    = "tcp"
    security_groups = [aws_security_group.apiserver-lb.id]
  }

  ingress {
    description = "Kubernetes API"
    from_port   = 6443
    to_port     = 6443
    protocol    = "tcp"
    self = true
  }

  ingress {
    description = "etcd"
    from_port   = 2379
    to_port     = 2379
    protocol    = "tcp"
    self = true
  }

  ingress {
    description = "etcd"
    from_port   = 2380
    to_port     = 2380
    protocol    = "tcp"
    self = true
  }

  ingress {
    description = "Node Port Services"
    from_port   = 30000
    to_port     = 32767
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    security_groups = [aws_security_group.bastion.id]
  }

  ingress {
    description = "Kubelet API"
    from_port   = 10250
    to_port     = 10250
    protocol    = "tcp"
    self = true
  }

  ingress {
    description = "Kubelet API"
    from_port   = 10250
    to_port     = 10250
    protocol    = "tcp"
    security_groups = [aws_security_group.node.id]
  }

  ingress {
    description = "bgp (calico)"
    from_port   = 179
    to_port     = 179
    protocol    = "tcp"
    self = true
  }

  ingress {
    description = "bgp (calico)"
    from_port   = 179
    to_port     = 179
    protocol    = "tcp"
    security_groups = [aws_security_group.node.id]
  }

  ingress {
    description = "IP-in-IP (calico)"
    protocol    = 4
  from_port = 0
  to_port = 0
    self = true
  }

  ingress {
    description = "IP-in-IP (calico)"
    protocol    = 4
  from_port = 0
  to_port = 0
    security_groups = [aws_security_group.node.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = join("-",[var.cluster_name,"controlplane"])
    (local.cluster_tag) = "shared"
    (local.role_tag) = "controlplane"
  }
}
