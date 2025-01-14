package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/calxus/tflint-ruleset-tofu/internal"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type TerraformVariableValidationRuleSeparationRule struct {
	tflint.DefaultRule
}

func NewTerraformVariableValidationRuleSeparationRule() *TerraformVariableValidationRuleSeparationRule {
	return &TerraformVariableValidationRuleSeparationRule{}
}

func (r *TerraformVariableValidationRuleSeparationRule) Name() string {
	return "terraform_variable_validation_rule_separation"
}

func (r *TerraformVariableValidationRuleSeparationRule) Enabled() bool {
	return true
}

func (r *TerraformVariableValidationRuleSeparationRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformVariableValidationRuleSeparationRule) Link() string {
	return ""
}

func (r *TerraformVariableValidationRuleSeparationRule) Check(runner tflint.Runner) error {
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
		attrs := internal.SortAttributesByStartLine(variable.Body.Attributes)

		if !checkVariableValidationRulesAreDividedByOneEmptyLine(variable, attrs) {
			if err := runner.EmitIssue(
				r,
				"Variable validation rules must be divided by one empty line",
				hcl.Range{
					Filename: variable.DefRange.Filename,
					Start:    variable.DefRange.Start,
					End:      variable.Body.Blocks[len(variable.Body.Blocks)-1].DefRange.End,
				},
			); err != nil {
				return err
			}
			continue
		}
	}

	return nil
}

func checkVariableValidationRulesAreDividedByOneEmptyLine(variable *hclext.Block, attrs []*hclext.Attribute) bool {
	l := attrs[len(attrs)-1].Range.End.Line
	for i := range variable.Body.Blocks {
		if variable.Body.Blocks[i].DefRange.Start.Line != l+2 {
			return false
		}

		vAttrs := internal.SortAttributesByStartLine(variable.Body.Blocks[i].Body.Attributes)
		l = vAttrs[len(vAttrs)-1].Range.End.Line + 1
	}

	return true
}
