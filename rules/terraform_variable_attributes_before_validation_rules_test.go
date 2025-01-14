package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/calxus/tflint-ruleset-tofu/internal"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformVariableAttributesBeforeValidationRules(t *testing.T) {
	tests := []internal.Test{
		{
			Name: "no validation before attributes",
			Content: `
variable "example" {
  description = "description"
  type = string
  default = "default"

  validation {
    error_message = "message"
    condition     = var.example
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "found validation before attributes",
			Content: `
variable "example" {
  validation {
    error_message = "message"
    condition     = var.example
  }
  description = "description"
  type = string
  default = "default"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformVariableAttributesBeforeValidationRulesRule(),
					Message: "Variable validation rules must be defined after attributes",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   2,
							Column: 19,
						},
					},
				},
			},
		},
	}

	internal.ExecuteTests(t, NewTerraformVariableAttributesBeforeValidationRulesRule(), tests)
}
