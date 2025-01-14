package rules

import (
	"fmt"

	"github.com/calxus/tflint-ruleset-tofu/internal"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type TerraformVariableAttributeOrderRule struct {
	tflint.DefaultRule
}

type terraformVariableAttributeOrderRuleConfig struct {
	Order []string `hclext:"order,optional"`
}

func NewTerraformVariableAttributeOrderRule() *TerraformVariableAttributeOrderRule {
	return &TerraformVariableAttributeOrderRule{}
}

func (r *TerraformVariableAttributeOrderRule) Name() string {
	return "terraform_variable_attribute_order"
}

func (r *TerraformVariableAttributeOrderRule) Enabled() bool {
	return true
}

func (r *TerraformVariableAttributeOrderRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformVariableAttributeOrderRule) Link() string {
	return ""
}

func (r *TerraformVariableAttributeOrderRule) config(runner tflint.Runner) (*terraformVariableAttributeOrderRuleConfig, error) {
	config := &terraformVariableAttributeOrderRuleConfig{}

	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return nil, err
	}

	dv := []string{"description", "type", "default"}
	if config.Order == nil {
		config.Order = dv
	}

	return config, nil
}

func (r *TerraformVariableAttributeOrderRule) Check(runner tflint.Runner) error {
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
					Attributes: []hclext.AttributeSchema{{Name: "description"}, {Name: "type"}, {Name: "default"}},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})

	if err != nil {
		return err
	}

	for _, variable := range body.Blocks {
		attrs := internal.SortAttributesByStartLine(variable.Body.Attributes)

		if !internal.CheckAttributeOrder(attrs, config.Order) {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("Variable attributes must be ordered by %v", config.Order),
				variable.DefRange,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
