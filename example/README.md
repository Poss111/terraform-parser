# Example Terraform Project

This is a comprehensive example that demonstrates all major Terraform components that the parser can analyze.

## Structure

```
example/
├── providers.tf       # Provider configurations (AWS, Random)
├── variables.tf       # Variable declarations with types and defaults
├── locals.tf         # Local values and computed configurations
├── data.tf           # Data sources for external information
├── modules.tf        # Module calls (VPC, Security Groups, EC2, S3)
├── resources.tf      # AWS resources (ALB, CloudWatch, IAM, Route53)
├── outputs.tf        # Output values
├── terraform.tfvars  # Default variable values
├── dev.tfvars        # Development environment values
└── staging.tfvars.json  # Staging environment values (JSON format)
```

## Components Included

### Providers
- AWS (with primary and secondary regions)
- Random provider
- Backend configuration for S3

### Variables
- String, number, boolean, list, and map types
- Validation rules
- Default values and descriptions

### Locals
- Naming conventions
- Common tags
- Network calculations
- Conditional logic

### Data Sources
- AWS AMI lookup
- Account and region information
- Availability zones
- Existing S3 buckets and IAM roles
- Conditional VPC lookup

### Modules
- Official AWS VPC module
- Security group module
- EC2 instance module
- S3 bucket module
- Local custom module reference

### Resources
- VPC and networking components
- Application Load Balancer with target groups
- CloudWatch logs and metric alarms
- SNS topics for notifications
- IAM roles and instance profiles
- Route53 health checks

### TfVars Files
- `terraform.tfvars` - Production configuration (HCL format)
- `dev.tfvars` - Development configuration (HCL format)
- `staging.tfvars.json` - Staging configuration (JSON format)

## Testing the Parser

Run the parser on this example:

```bash
./terraform-parser example
```

With verbose output:
```bash
./terraform-parser example -v
```

Save to file:
```bash
./terraform-parser example -o example-output.json
```

