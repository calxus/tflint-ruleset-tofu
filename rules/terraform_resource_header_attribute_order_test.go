package rules

import (
	"testing"

	"github.com/calxus/tflint-ruleset-tofu/internal"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/hashicorp/hcl/v2"
)

func Test_TerraformResourceHeaderAttributeOrder(t *testing.T) {
	tests := []internal.Test{
		{
			Name: "no meta arguments",
			Content: `
resource "aws_instance" "example" {
  ami           = "ami-a1b2c3d4"
  instance_type = "t2.micro"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "provider meta argument correct",
			Content: `
resource "aws_instance" "example" {
  provider      = aws.default
  ami           = "ami-a1b2c3d4"
  instance_type = "t2.micro"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "count meta argument correct",
			Content: `
resource "aws_instance" "example" {
  count         = 1
  ami           = "ami-a1b2c3d4"
  instance_type = "t2.micro"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "for_each meta argument correct",
			Content: `
resource "aws_instance" "example" {
  for_each      = local.instances
  ami           = "ami-a1b2c3d4"
  instance_type = "t2.micro"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "multiple meta argument correct",
			Content: `
resource "aws_instance" "example" {
  provider      = aws.default
  count         = 1
  ami           = "ami-a1b2c3d4"
  instance_type = "t2.micro"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "provider meta argument incorrect",
			Content: `
resource "aws_instance" "example" {
  ami           = "ami-a1b2c3d4"
  provider      = aws.default
  instance_type = "t2.micro"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceHeaderAttributeOrderRule(),
					Message: "The provider meta-argument must be the first attribute in resource definition",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 3,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 30,
						},
					},
				},},
		},
		{
			Name: "count meta argument incorrect",
			Content: `
resource "aws_instance" "example" {
  ami           = "ami-a1b2c3d4"
  count         = 1
  instance_type = "t2.micro"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceHeaderAttributeOrderRule(),
					Message: "The count meta-argument must be defined at the top of the resource definition",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 3,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 20,
						},
					},
				},},
		},
		{
			Name: "for_each meta argument incorrect",
			Content: `
resource "aws_instance" "example" {
  ami           = "ami-a1b2c3d4"
  for_each      = local.instances
  instance_type = "t2.micro"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceHeaderAttributeOrderRule(),
					Message: "The for_each meta-argument must be defined at the top of the resource definition",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 3,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 34,
						},
					},
				},},
		},
		{
			Name: "multiple meta argument incorrect",
			Content: `
resource "aws_instance" "example" {
  count         = 1
  provider      = aws.default
  ami           = "ami-a1b2c3d4"
  instance_type = "t2.micro"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformResourceHeaderAttributeOrderRule(),
					Message: "The provider meta-argument must be the first attribute in resource definition",
					Range: hcl.Range{
						Filename: "resource.tf",
						Start: hcl.Pos{
							Line:   4,
							Column: 3,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 30,
						},
					},
				},},
		},
	}

	internal.ExecuteTests(t, NewTerraformResourceHeaderAttributeOrderRule(), tests)
}
