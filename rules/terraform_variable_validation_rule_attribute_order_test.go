package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/calxus/tflint-ruleset-tofu/internal"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformVariableValidationRuleAttributeOrder(t *testing.T) {
	tests := []internal.Test{
		{
			Name: "no order issue",
			Content: `
variable "example" {
  validation {
    error_message = "message"
    condition     = var.example
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "found order issue",
			Content: `
variable "example" {
  validation {
    condition     = var.example
    error_message = "message"
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformVariableValidationRuleAttributeOrderRule(),
					Message: "Variable validation rule attributes must be ordered by [error_message condition]",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 5,
						},
						End: hcl.Pos{
							Line:   5,
							Column: 30,
						},
					},
				},
			},
		},
	}

	internal.ExecuteTests(t, NewTerraformVariableValidationRuleAttributeOrderRule(), tests)
}
