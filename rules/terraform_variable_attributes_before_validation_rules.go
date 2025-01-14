package rules

import (
	"github.com/calxus/tflint-ruleset-tofu/internal"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type TerraformVariableAttributesBeforeValidationRulesRule struct {
	tflint.DefaultRule
}

func NewTerraformVariableAttributesBeforeValidationRulesRule() *TerraformVariableAttributesBeforeValidationRulesRule {
	return &TerraformVariableAttributesBeforeValidationRulesRule{}
}

func (r *TerraformVariableAttributesBeforeValidationRulesRule) Name() string {
	return "terraform_variable_attributes_before_validation_rules"
}

func (r *TerraformVariableAttributesBeforeValidationRulesRule) Enabled() bool {
	return true
}

func (r *TerraformVariableAttributesBeforeValidationRulesRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformVariableAttributesBeforeValidationRulesRule) Link() string {
	return ""
}

func (r *TerraformVariableAttributesBeforeValidationRulesRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		// This rule does not evaluate child modules.
		return nil
	}

	body, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{{Name: "description"}, {Name: "type"}, {Name: "default"}},
					Blocks:     []hclext.BlockSchema{{Type: "validation"}},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})

	if err != nil {
		return err
	}

	for _, variable := range body.Blocks {
		attrs := internal.SortAttributesByStartLine(variable.Body.Attributes)

		if !checkValidationRulesAreDefinedAfterAttributes(variable, attrs) {
			if err := runner.EmitIssue(
				r,
				"Variable validation rules must be defined after attributes",
				variable.DefRange,
			); err != nil {
				return err
			}
			continue
		}
	}

	return nil
}

func checkValidationRulesAreDefinedAfterAttributes(variable *hclext.Block, attrs []*hclext.Attribute) bool {
	return !(len(variable.Body.Blocks) > 0 && attrs[0].Range.Start.Line > variable.Body.Blocks[0].DefRange.Start.Line)
}
