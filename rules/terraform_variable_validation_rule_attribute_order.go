package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/calxus/tflint-ruleset-tofu/internal"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type TerraformVariableValidationRuleAttributeOrderRule struct {
	tflint.DefaultRule
}

type terraformVariableValidationRuleAttributeOrderRuleConfig struct {
	Order []string `hclext:"order,optional"`
}

func NewTerraformVariableValidationRuleAttributeOrderRule() *TerraformVariableValidationRuleAttributeOrderRule {
	return &TerraformVariableValidationRuleAttributeOrderRule{}
}

func (r *TerraformVariableValidationRuleAttributeOrderRule) Name() string {
	return "terraform_variable_validation_rule_attribute_order"
}

func (r *TerraformVariableValidationRuleAttributeOrderRule) Enabled() bool {
	return true
}

func (r *TerraformVariableValidationRuleAttributeOrderRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformVariableValidationRuleAttributeOrderRule) Link() string {
	return ""
}

func (r *TerraformVariableValidationRuleAttributeOrderRule) config(runner tflint.Runner) (*terraformVariableAttributeOrderRuleConfig, error) {
	config := &terraformVariableAttributeOrderRuleConfig{}

	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return nil, err
	}

	dv := []string{"error_message", "condition"}
	if config.Order == nil {
		config.Order = dv
	}

	return config, nil
}

func (r *TerraformVariableValidationRuleAttributeOrderRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	config, err := r.config(runner)
	if err != nil {
		return fmt.Errorf("failed to parse rule config: %w", err)
	}

	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Blocks: []hclext.BlockSchema{{Type: "validation", Body: &hclext.BodySchema{
						Attributes: []hclext.AttributeSchema{{Name: "condition"}, {Name: "error_message"}},
					}}},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})

	if err != nil {
		return err
	}

	for _, variable := range body.Blocks {
		for _, validation := range variable.Body.Blocks {
			attrs := internal.SortAttributesByStartLine(validation.Body.Attributes)

			if !internal.CheckAttributeOrder(attrs, config.Order) {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Variable validation rule attributes must be ordered by %v", config.Order),
					hcl.Range{
						Filename: attrs[0].Range.Filename,
						Start:    attrs[0].Range.Start,
						End:      attrs[len(attrs)-1].Range.End,
					},
				); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
