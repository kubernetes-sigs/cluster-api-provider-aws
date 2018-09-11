resource "aws_security_group" "cluster_default" {
  name        = "${format("%s-default", var.cluster_name)}"
  description = "${format("%s default security policy", var.cluster_name)}"
  vpc_id      = "${module.vpc.vpc_id}"

  tags {
    Name = "${var.cluster_name}"
  }
}

resource "aws_security_group_rule" "deault_egress" {
  type              = "egress"
  security_group_id = "${aws_security_group.cluster_default.id}"

  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "default_ingress_ssh" {
  type              = "ingress"
  security_group_id = "${aws_security_group.cluster_default.id}"

  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 22
  to_port     = 22
}
