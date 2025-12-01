# Default values for all environments
aws_region     = "us-east-1"
secondary_region = "us-west-2"
project_name   = "my-app"
environment    = "production"

# Network configuration
vpc_cidr = "10.0.0.0/16"
availability_zones = [
  "us-east-1a",
  "us-east-1b",
  "us-east-1c"
]

# NAT Gateway
enable_nat_gateway = true

# EC2 configuration
instance_type  = "t3.small"
instance_count = 3

# Monitoring
enable_monitoring = true

# Additional tags
tags = {
  Team       = "DevOps"
  CostCenter = "Engineering"
  Compliance = "PCI-DSS"
}

complex_json = {
  users = [
    {
      id = 1
      name = "Alice"
      roles = ["admin", "developer"]
      settings = {
        theme      = "dark"
        languages  = ["en", "es"]
        marketing  = { subscribed = true, frequency = "weekly" }
      }
    },
    {
      id = 2
      name = "Bob"
      roles = ["viewer"]
      settings = {
        theme      = "light"
        languages  = ["en"]
        marketing  = { subscribed = false, frequency = "" }
      }
    }
  ]
  features = {
    authentication = {
      enabled    = true
      providers  = ["saml", "google", "github"]
      config     = {
        saml = { enabled = true, entity_id = "saml-entity-001" }
        google = { enabled = true, client_id = "google-client-xyz" }
        github = { enabled = false }
      }
    }
    logging = {
      enabled = true
      level   = "info"
      destinations = {
        cloudwatch = { stream = "main-logs", retention_days = 30 }
        s3         = { bucket = "my-app-logs", prefix = "2024/", encryption = true }
      }
    }
  }
  metadata = {
    created_by = "terraform"
    timestamp  = "2024-06-13T12:34:56Z"
    tags       = ["production", "v2", "json"]
  }
}
