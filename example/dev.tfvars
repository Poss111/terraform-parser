# Development environment configuration
aws_region     = "us-east-1"
project_name   = "my-app"
environment    = "dev"

# Network - smaller CIDR for dev
vpc_cidr = "10.10.0.0/16"
availability_zones = ["us-east-1a"]

# No NAT Gateway in dev to save costs
enable_nat_gateway = false

# EC2 - minimal resources for dev
instance_type  = "t3.micro"
instance_count = 1

# No detailed monitoring in dev
enable_monitoring = false

# Dev-specific tags
tags = {
  Team        = "Development"
  AutoShutdown = "true"
  Owner       = "dev-team"
}

