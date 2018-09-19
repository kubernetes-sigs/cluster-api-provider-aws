resource "aws_security_group" "cluster_default" {
  name        = "${format("%s-default", var.environment_id)}"
  description = "${format("%s default security policy", var.environment_id)}"
  vpc_id      = "${module.vpc.vpc_id}"

  tags {
    Name = "${var.environment_id}-worker-sg"
  }
}

resource "aws_security_group_rule" "default_egress" {
  type              = "egress"
  security_group_id = "${aws_security_group.cluster_default.id}"

  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "default_ingress" {
  type              = "ingress"
  security_group_id = "${aws_security_group.cluster_default.id}"

  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "allow_all" {
  type              = "ingress"
  security_group_id = "${aws_security_group.cluster_default.id}"

  cidr_blocks = ["0.0.0.0/0"]
  from_port       = 22
  to_port         = 22
  description     = "SSH"
  protocol        = "tcp"
}
