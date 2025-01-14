package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/calxus/tflint-ruleset-tofu/internal"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformVariableSeparation(t *testing.T) {
	tests := []internal.Test{
		{
			Name: "no attribute separator issue",
			Content: `
variable "example" {
  description = "description"
  type = object({
	key = string					
  })
  default = "default"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "found attribute separator issue",
			Content: `
variable "example" {
  description = "description"
  
  type = object({
	key = string					
  })

  default = "default"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformVariableAttributeSeparationRule(),
					Message: "No empty line is expected between attributes",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start: hcl.Pos{
							Line:   3,
							Column: 3,
						},
						End: hcl.Pos{
							Line:   9,
							Column: 22,
						},
					},
				},
			},
		},
	}

	internal.ExecuteTests(t, NewTerraformVariableAttributeSeparationRule(), tests)
}
