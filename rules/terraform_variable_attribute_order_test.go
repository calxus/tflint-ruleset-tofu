package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/calxus/tflint-ruleset-tofu/internal"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformVariableAttributeOrder(t *testing.T) {
	tests := []internal.Test{
		{
			Name: "no order issue",
			Content: `
variable "example" {
  description = "description"
  type = string
  default = "default"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no order issue no default",
			Content: `
variable "example" {
  description = "description"
  type = string
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no order issue no type",
			Content: `
variable "example" {
  description = "description"
  default = "default"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "found order issue",
			Content: `
variable "example" {
  type = string
  default = "default"
  description = "description"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformVariableAttributeOrderRule(),
					Message: "Variable attributes must be ordered by [description type default]",
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
		{
			Name: "found order issue multiline default",
			Content: `
variable "example" {
  description = "description"
  default = [
	"foo",
    "bar",
  ]
  type = list(string)
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformVariableAttributeOrderRule(),
					Message: "Variable attributes must be ordered by [description type default]",
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

	internal.ExecuteTests(t, NewTerraformVariableAttributeOrderRule(), tests)
}
