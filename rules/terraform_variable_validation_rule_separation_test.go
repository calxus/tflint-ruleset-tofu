package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/calxus/tflint-ruleset-tofu/internal"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformVariableValidationRuleSeparation(t *testing.T) {
	tests := []internal.Test{
		{
			Name: "no validation separation issue",
			Content: `
variable "example" {
  description = "description"
  type = string
  default = "default"

  validation {
    error_message = "message"
    condition     = var.example
  }

  validation {
    error_message = "message2"
    condition     = var.example
  }
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "found validation separation issue",
			Content: `
variable "example" {
  description = "description"
  type = string
  default = "default"

  validation {
    error_message = "message"
    condition     = var.example
  }
  validation {
    error_message = "message2"
    condition     = var.example
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformVariableValidationRuleSeparationRule(),
					Message: "Variable validation rules must be divided by one empty line",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   11,
							Column: 13,
						},
					},
				},
			},
		},
	}

	internal.ExecuteTests(t, NewTerraformVariableValidationRuleSeparationRule(), tests)
}
