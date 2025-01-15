package rules

import (
	"fmt"
	"sort"
	"github.com/hashicorp/hcl/v2"
	// "github.com/calxus/tflint-ruleset-tofu/internal"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	// "github.com/terraform-linters/tflint-plugin-sdk/hclext"
)

type TerraformResourceHeaderAttributeOrderRule struct {
	tflint.DefaultRule
}

type terraformResourceHeaderAttributeOrderRuleConfig struct {
	Order []string `hclext:"order,optional"`
}

func NewTerraformResourceHeaderAttributeOrderRule() *TerraformResourceHeaderAttributeOrderRule {
	return &TerraformResourceHeaderAttributeOrderRule{}
}

func (r *TerraformResourceHeaderAttributeOrderRule) Name() string {
	return "terraform_resource_header_attribute_order"
}

func (r *TerraformResourceHeaderAttributeOrderRule) Enabled() bool {
	return true
}

func (r *TerraformResourceHeaderAttributeOrderRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformResourceHeaderAttributeOrderRule) Link() string {
	return ""
}

func (r *TerraformResourceHeaderAttributeOrderRule) config(runner tflint.Runner) (*terraformResourceHeaderAttributeOrderRuleConfig, error) {
	config := &terraformResourceHeaderAttributeOrderRuleConfig{}

	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return nil, err
	}

	dv := []string{"provider", "for_each", "count"}
	if config.Order == nil {
		config.Order = dv
	}

	return config, nil
}

func (r *TerraformResourceHeaderAttributeOrderRule) Check(runner tflint.Runner) error {
	files, _ := runner.GetFiles()
	for _, f := range files {
		content, _, _ := f.Body.PartialContent(&hcl.BodySchema{ 
			Blocks: []hcl.BlockHeaderSchema{ 
				{ 
					Type:       "resource", 
					LabelNames: []string{"type", "name"}, 
				}, 
			}, 
		})
		for _, resource := range content.Blocks {
			attrs, _ := resource.Body.JustAttributes()
			sorted_attrs := SortAttributesByStartLine(attrs)
			if (AttributesContains(sorted_attrs, "provider")) {
				if (sorted_attrs[0].Name != "provider") {
					runner.EmitIssue(
						r,
						fmt.Sprintf("The provider meta-argument must be the first attribute in resource definition if present"),
						sorted_attrs[0].Range,
					);
					continue
				}
				if (AttributesContains(sorted_attrs, "count") || AttributesContains(sorted_attrs, "for_each")) {
					if !((sorted_attrs[1].Name == "count") || (sorted_attrs[1].Name == "for_each")) {
						runner.EmitIssue(
							r,
							fmt.Sprintf("The count/for_each meta-arguments must appear after the provider meta-argument if present"),
							sorted_attrs[1].Range,
						);
						continue
					}
				}
			} else if (AttributesContains(sorted_attrs, "count") || AttributesContains(sorted_attrs, "for_each")) {
				if !((sorted_attrs[0].Name == "count") || (sorted_attrs[0].Name == "for_each")) {
					runner.EmitIssue(
						r,
						fmt.Sprintf("The count/for_each meta-arguments must be defined at the top of the resource definition"),
						sorted_attrs[1].Range,
					);
					continue
				}
			}
		}
	}
	return nil
}

func SortAttributesByStartLine(attributes hcl.Attributes) []*hcl.Attribute {
	var attrs []*hcl.Attribute
	for _, a := range attributes {
		attrs = append(attrs, a)
	}
	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].Range.Start.Line < attrs[j].Range.Start.Line
	})

	return attrs
}

func AttributesContains(attributes []*hcl.Attribute, attribute string) bool {
	for _, attr := range attributes {
		if (attr.Name == attribute) {
			return true
		}
	}
	return false
}
