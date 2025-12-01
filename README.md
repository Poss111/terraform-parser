# Terraform Parser

A Go CLI tool that parses Terraform code and generates a JSON breakdown of resources, modules, and providers.

Built with [Cobra CLI](https://github.com/spf13/cobra) for a modern command-line experience.

## Installation

```bash
go mod download
go build -o terraform-parser
```

## Usage

```bash
terraform-parser <terraform-directory> [flags]
```

### Flags

- `-o, --output <file>` - Write output to file instead of stdout
- `-p, --pretty` - Pretty print JSON output (default: true)
- `-v, --verbose` - Verbose output (show warnings and parsing progress)
- `-h, --help` - Show help message

### Examples

Parse a Terraform directory:
```bash
./terraform-parser ./my-terraform-project
```

Save output to a file:
```bash
./terraform-parser ./infra -o output.json
```

Compact JSON output:
```bash
./terraform-parser ./infra -p=false
```

Verbose mode with file output:
```bash
./terraform-parser ./infra -v -o result.json
```

## Example Output

```json
{
  "resources": [
    {
      "type": "aws_instance",
      "name": "web",
      "file": "/path/to/main.tf",
      "attributes": {
        "ami": "ami-12345678",
        "instance_type": "t2.micro"
      }
    }
  ],
  "modules": [
    {
      "name": "vpc",
      "source": "terraform-aws-modules/vpc/aws",
      "file": "/path/to/main.tf",
      "attributes": {
        "source": "terraform-aws-modules/vpc/aws",
        "version": "3.0.0",
        "cidr": "10.0.0.0/16"
      }
    }
  ],
  "providers": [
    {
      "name": "aws",
      "file": "/path/to/main.tf",
      "attributes": {
        "region": "us-west-2"
      }
    }
  ]
}
```

## Features

- 🔍 **Recursive Scanning** - Automatically finds all `.tf` files in subdirectories
- 📦 **Resource Extraction** - Captures all resource types, names, and attributes
- 🔗 **Module Detection** - Identifies module calls with source information
- ⚙️ **Provider Discovery** - Detects provider configurations and requirements
- 📄 **JSON Output** - Structured JSON for easy processing and integration
- 💾 **File Output** - Save results directly to a file
- 🎨 **Flexible Formatting** - Pretty or compact JSON output
- 📢 **Verbose Mode** - See detailed parsing progress

## Technical Details

- Uses the official HashiCorp HCL parser (`github.com/hashicorp/hcl/v2`)
- Built with Cobra CLI framework for professional command-line interface
- Handles complex attribute values (strings, numbers, booleans, lists, maps)
- Provides helpful error messages and validation
- Robust error handling for malformed files


