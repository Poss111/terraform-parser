package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/spf13/cobra"
	"github.com/zclconf/go-cty/cty"
)

// TerraformBreakdown represents the complete analysis of a Terraform codebase
type TerraformBreakdown struct {
	Resources []Resource        `json:"resources"`
	Modules   []Module          `json:"modules"`
	Providers []Provider        `json:"providers"`
	Variables []Variable        `json:"variables"`
	TfVars    map[string]TfVars `json:"tfvars"`
}

// Resource represents a Terraform resource
type Resource struct {
	Type       string            `json:"type"`
	Name       string            `json:"name"`
	File       string            `json:"file"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// Module represents a Terraform module call
type Module struct {
	Name       string            `json:"name"`
	Source     string            `json:"source"`
	File       string            `json:"file"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// Provider represents a Terraform provider configuration
type Provider struct {
	Name       string            `json:"name"`
	Alias      string            `json:"alias,omitempty"`
	File       string            `json:"file"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// Variable represents a Terraform variable declaration
type Variable struct {
	Name        string `json:"name"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Default     string `json:"default,omitempty"`
	File        string `json:"file"`
}

// TfVars represents variable values from a .tfvars file
type TfVars struct {
	File   string            `json:"file"`
	Values map[string]string `json:"values"`
}

var (
	outputFile  string
	prettyPrint bool
	verbose     bool
)

var rootCmd = &cobra.Command{
	Use:   "terraform-parser <directory>",
	Short: "Parse Terraform code and output JSON breakdown",
	Long: `Terraform Parser analyzes Terraform code and generates a structured JSON breakdown
of resources, modules, and providers found in .tf files.

The tool recursively scans the specified directory for Terraform files and extracts:
  - Resources: All resource blocks with their types, names, and attributes
  - Modules: Module calls with source information and parameters
  - Providers: Provider configurations and requirements
  - Variables: Variable declarations from .tf files
  - TfVars: Variable values from .tfvars and .tfvars.json files`,
	Example: `  # Parse a Terraform directory
  terraform-parser ./my-terraform-project

  # Save output to a file
  terraform-parser ./infra -o output.json

  # Compact JSON output
  terraform-parser ./infra -p=false`,
	Args: cobra.ExactArgs(1),
	RunE: runParse,
}

func init() {
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write output to file instead of stdout")
	rootCmd.Flags().BoolVarP(&prettyPrint, "pretty", "p", true, "Pretty print JSON output")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output (show warnings)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runParse(cmd *cobra.Command, args []string) error {
	dirPath := args[0]

	// Verify directory exists
	if info, err := os.Stat(dirPath); err != nil || !info.IsDir() {
		return fmt.Errorf("%s is not a valid directory", dirPath)
	}

	breakdown, err := parseTerraformDirectory(dirPath)
	if err != nil {
		return fmt.Errorf("error parsing Terraform directory: %w", err)
	}

	// Marshal JSON based on pretty print flag
	var output []byte
	if prettyPrint {
		output, err = json.MarshalIndent(breakdown, "", "  ")
	} else {
		output, err = json.Marshal(breakdown)
	}
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	// Output to file or stdout
	if outputFile != "" {
		if err := os.WriteFile(outputFile, output, 0644); err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "Output written to %s\n", outputFile)
		}
	} else {
		fmt.Println(string(output))
	}

	return nil
}

func parseTerraformDirectory(dirPath string) (*TerraformBreakdown, error) {
	breakdown := &TerraformBreakdown{
		Resources: []Resource{},
		Modules:   []Module{},
		Providers: []Provider{},
		Variables: []Variable{},
		TfVars:    make(map[string]TfVars),
	}

	parser := hclparse.NewParser()

	// Walk through the directory and find all .tf and .tfvars files
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		baseName := filepath.Base(path)

		// Handle .tfvars files (both .tfvars and .tfvars.json)
		if ext == ".tfvars" || strings.HasSuffix(baseName, ".tfvars.json") {
			if err := parseTfVarsFile(parser, path, breakdown); err != nil {
				if verbose {
					fmt.Fprintf(os.Stderr, "Warning: Error parsing %s: %v\n", path, err)
				}
			} else if verbose {
				fmt.Fprintf(os.Stderr, "Parsed tfvars: %s\n", path)
			}
			return nil
		}

		// Handle .tf files
		if ext != ".tf" {
			return nil
		}

		// Parse the file
		file, diags := parser.ParseHCLFile(path)
		if diags.HasErrors() {
			if verbose {
				fmt.Fprintf(os.Stderr, "Warning: Error parsing %s: %v\n", path, diags)
			}
			return nil // Continue processing other files
		}

		if verbose {
			fmt.Fprintf(os.Stderr, "Parsed: %s\n", path)
		}

		// Extract resources, modules, providers, and variables from the file
		extractFromFile(file.Body, path, breakdown)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return breakdown, nil
}

func extractFromFile(body hcl.Body, filePath string, breakdown *TerraformBreakdown) {
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "resource", LabelNames: []string{"type", "name"}},
			{Type: "module", LabelNames: []string{"name"}},
			{Type: "provider", LabelNames: []string{"name"}},
			{Type: "terraform"},
			{Type: "variable", LabelNames: []string{"name"}},
			{Type: "output"},
			{Type: "data", LabelNames: []string{"type", "name"}},
		},
	}

	content, _, diags := body.PartialContent(schema)
	if diags.HasErrors() {
		return
	}

	for _, block := range content.Blocks {
		switch block.Type {
		case "resource":
			if len(block.Labels) >= 2 {
				resource := Resource{
					Type:       block.Labels[0],
					Name:       block.Labels[1],
					File:       filePath,
					Attributes: extractAttributes(block.Body),
				}
				breakdown.Resources = append(breakdown.Resources, resource)
			}

		case "module":
			if len(block.Labels) >= 1 {
				module := Module{
					Name:       block.Labels[0],
					File:       filePath,
					Attributes: extractAttributes(block.Body),
				}
				// Extract source attribute specifically
				if source, ok := module.Attributes["source"]; ok {
					module.Source = source
				}
				breakdown.Modules = append(breakdown.Modules, module)
			}

		case "provider":
			if len(block.Labels) >= 1 {
				provider := Provider{
					Name:       block.Labels[0],
					File:       filePath,
					Attributes: extractAttributes(block.Body),
				}
				// Extract alias if present
				if alias, ok := provider.Attributes["alias"]; ok {
					provider.Alias = alias
				}
				breakdown.Providers = append(breakdown.Providers, provider)
			}

		case "variable":
			if len(block.Labels) >= 1 {
				variable := Variable{
					Name: block.Labels[0],
					File: filePath,
				}
				attrs := extractAttributes(block.Body)
				if varType, ok := attrs["type"]; ok {
					variable.Type = varType
				}
				if desc, ok := attrs["description"]; ok {
					variable.Description = desc
				}
				if def, ok := attrs["default"]; ok {
					variable.Default = def
				}
				breakdown.Variables = append(breakdown.Variables, variable)
			}

		case "terraform":
			// Extract required_providers from terraform block
			extractRequiredProviders(block.Body, filePath, breakdown)
		}
	}
}

func extractAttributes(body hcl.Body) map[string]string {
	attrs := make(map[string]string)
	
	bodyAttrs, diags := body.JustAttributes()
	if diags.HasErrors() {
		return attrs
	}

	for name, attr := range bodyAttrs {
		// Try to evaluate as a literal value first
		val, diags := attr.Expr.Value(nil)
		if !diags.HasErrors() {
			formatted := formatCtyValue(val)
			if formatted != "" {
				attrs[name] = formatted
			}
		}
	}

	return attrs
}

// formatCtyValue converts a cty.Value to a readable string
func formatCtyValue(val cty.Value) string {
	if val.IsNull() {
		return "null"
	}

	switch val.Type() {
	case cty.String:
		return val.AsString()
	case cty.Number:
		bf := val.AsBigFloat()
		if i, acc := bf.Int64(); acc == 0 {
			return fmt.Sprintf("%d", i)
		}
		f, _ := bf.Float64()
		return fmt.Sprintf("%g", f)
	case cty.Bool:
		if val.True() {
			return "true"
		}
		return "false"
	default:
		// For complex types like lists, maps, objects
		if val.Type().IsListType() || val.Type().IsTupleType() {
			var items []string
			it := val.ElementIterator()
			for it.Next() {
				_, v := it.Element()
				items = append(items, formatCtyValue(v))
			}
			return "[" + strings.Join(items, ", ") + "]"
		}
		if val.Type().IsMapType() || val.Type().IsObjectType() {
			var pairs []string
			it := val.ElementIterator()
			for it.Next() {
				k, v := it.Element()
				pairs = append(pairs, fmt.Sprintf("%s: %s", formatCtyValue(k), formatCtyValue(v)))
			}
			return "{" + strings.Join(pairs, ", ") + "}"
		}
		// Fallback to JSON
		data, err := json.Marshal(val)
		if err == nil {
			return string(data)
		}
		return fmt.Sprintf("%v", val)
	}
}

func extractRequiredProviders(body hcl.Body, filePath string, breakdown *TerraformBreakdown) {
	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "required_providers"},
		},
	}

	content, _, diags := body.PartialContent(schema)
	if diags.HasErrors() {
		return
	}

	for _, block := range content.Blocks {
		if block.Type == "required_providers" {
			attrs := extractAttributes(block.Body)
			for providerName := range attrs {
				// Check if this provider is already in the list
				found := false
				for _, p := range breakdown.Providers {
					if p.Name == providerName && p.File == filePath {
						found = true
						break
					}
				}
				if !found {
					breakdown.Providers = append(breakdown.Providers, Provider{
						Name: providerName,
						File: filePath,
					})
				}
			}
		}
	}
}

// parseTfVarsFile parses a .tfvars or .tfvars.json file
func parseTfVarsFile(parser *hclparse.Parser, filePath string, breakdown *TerraformBreakdown) error {
	// Determine if it's JSON or HCL
	var file *hcl.File
	var diags hcl.Diagnostics

	if strings.HasSuffix(filePath, ".json") {
		file, diags = parser.ParseJSONFile(filePath)
	} else {
		file, diags = parser.ParseHCLFile(filePath)
	}

	if diags.HasErrors() {
		return fmt.Errorf("parse error: %v", diags)
	}

	// Extract all attributes as variable values
	attrs, diags := file.Body.JustAttributes()
	if diags.HasErrors() {
		return fmt.Errorf("attribute extraction error: %v", diags)
	}

	values := make(map[string]string)
	for name, attr := range attrs {
		val, valDiags := attr.Expr.Value(nil)
		if !valDiags.HasErrors() {
			values[name] = formatCtyValue(val)
		}
	}

	// Use the base name without extension as the key
	baseName := filepath.Base(filePath)
	key := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	if strings.HasSuffix(key, ".tfvars") {
		key = strings.TrimSuffix(key, ".tfvars")
	}
	if key == "" {
		key = "terraform"
	}

	breakdown.TfVars[key] = TfVars{
		File:   filePath,
		Values: values,
	}

	return nil
}

