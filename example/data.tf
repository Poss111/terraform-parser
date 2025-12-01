# Data source: Get latest Amazon Linux 2 AMI
data "aws_ami" "amazon_linux_2" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

# Data source: Get current AWS account ID
data "aws_caller_identity" "current" {}

# Data source: Get current AWS region
data "aws_region" "current" {}

# Data source: Get available availability zones
data "aws_availability_zones" "available" {
  state = "available"
}

# Data source: Get existing S3 bucket for logging
data "aws_s3_bucket" "logs" {
  bucket = "${var.project_name}-logs-${data.aws_caller_identity.current.account_id}"
}

# Data source: Get existing IAM role
data "aws_iam_role" "ec2_role" {
  name = "EC2-DefaultRole"
}

# Data source: Get VPC information (if using existing VPC)
data "aws_vpc" "existing" {
  count = var.environment == "dev" ? 1 : 0
  
  tags = {
    Name = "dev-vpc"
  }
}

