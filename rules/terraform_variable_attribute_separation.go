package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/calxus/tflint-ruleset-tofu/internal"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type TerraformVariableAttributeSeparationRule struct {
	tflint.DefaultRule
}

func NewTerraformVariableAttributeSeparationRule() *TerraformVariableAttributeSeparationRule {
	return &TerraformVariableAttributeSeparationRule{}
}

func (r *TerraformVariableAttributeSeparationRule) Name() string {
	return "terraform_variable_attribute_separation"
}

func (r *TerraformVariableAttributeSeparationRule) Enabled() bool {
	return true
}

func (r *TerraformVariableAttributeSeparationRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformVariableAttributeSeparationRule) Link() string {
	return ""
}

func (r *TerraformVariableAttributeSeparationRule) Check(runner tflint.Runner) error {
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
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})

	if err != nil {
		return err
	}

	for _, variable := range body.Blocks {
		attrs := internal.SortAttributesByStartLine(variable.Body.Attributes)

		if !checkNoEmptyLineBetweenAttributes(variable, attrs) {
			if err := runner.EmitIssue(
				r,
				"No empty line is expected between attributes",
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

	return nil
}

func checkNoEmptyLineBetweenAttributes(variable *hclext.Block, attrs []*hclext.Attribute) bool {
	l := variable.DefRange.Start.Line + 1
	for i := range attrs {
		if attrs[i].Range.Start.Line != l {
			return false
		}
		l = attrs[i].Range.End.Line + 1
	}

	return true
}
