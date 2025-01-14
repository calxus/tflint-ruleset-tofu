package internal

import (
	"sort"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
)

func SortAttributesByStartLine(attributes hclext.Attributes) []*hclext.Attribute {
	var attrs []*hclext.Attribute
	for _, a := range attributes {
		attrs = append(attrs, a)
	}
	sort.Slice(attrs, func(i, j int) bool {
		return attrs[i].Range.Start.Line < attrs[j].Range.Start.Line
	})

	return attrs
}

func CheckAttributeOrder(attrs []*hclext.Attribute, order []string) bool {
	var o []string
	var a []*hclext.Attribute

	// remove non-existing attributes
	for i := range attrs {
		for j := range order {
			if attrs[i].Name == order[j] {
				o = append(o, order[j])
				a = append(a, attrs[i])
			}
		}
	}

	for i := range a {
		for j := range o {
			// if attribute exists in order slice
			if attrs[i].Name == order[j] {
				// check its position
				if attrs[i].Name != order[i] {
					// attribute is not at the expected position
					return false
				}
			}
		}
	}

	return true
}
