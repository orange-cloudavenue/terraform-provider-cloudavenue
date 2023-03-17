package superschema

import (
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"golang.org/x/exp/slices"
)

const (
	Required Param = iota
	Optional
	OptionalComputed
	Computed
)

type Param int

type Schema struct {
	MarkdownDescription string
	DeprecationMessage  string
	Attributes          map[string]schemaR.Attribute
}

func (s Schema) GetResource() schemaR.Schema {
	return schemaR.Schema{
		MarkdownDescription: s.MarkdownDescription,
		DeprecationMessage:  s.DeprecationMessage,
		Attributes:          s.Attributes,
	}
}

func (s Schema) GetDataSource() schemaD.Schema {
	sD := schemaD.Schema{
		MarkdownDescription: s.MarkdownDescription,
		DeprecationMessage:  s.DeprecationMessage,
		Attributes:          map[string]schemaD.Attribute{},
	}

	for k, v := range s.GetResource().GetAttributes() {
		sD.Attributes[k] = v
	}

	return sD
}

// SetParam will set the param for the schema for attributes in slices attributes.
// If you want to set the param for all attributes, you doesn't need to pass any attributes.
func (s Schema) SetParam(p Param, attributes ...string) Schema {
	for e, v := range s.Attributes {
		// the default params values
		required := false
		optional := false
		computed := false

		// if current attribute name is in attributes slice
		// or if attributes slice is empty, set the param.
		if slices.Contains(attributes, e) || len(attributes) == 0 {
			switch p {
			case Required:
				required = true
			case Optional:
				optional = true
			case OptionalComputed:
				optional = true
				computed = true
			case Computed:
				computed = true
			}
		}

		switch v.(type) {
		case schemaR.StringAttribute:
			attr := s.Attributes[e].(schemaR.StringAttribute)
			attr.Required = required
			attr.Computed = computed
			attr.Optional = optional
			s.Attributes[e] = attr

		case schemaR.NumberAttribute:
			attr := s.Attributes[e].(schemaR.NumberAttribute)
			attr.Required = required
			attr.Computed = computed
			attr.Optional = optional
			s.Attributes[e] = attr

		case schemaR.Int64Attribute:
			attr := s.Attributes[e].(schemaR.Int64Attribute)
			attr.Required = required
			attr.Computed = computed
			attr.Optional = optional
			s.Attributes[e] = attr

		case schemaR.BoolAttribute:
			attr := s.Attributes[e].(schemaR.BoolAttribute)
			attr.Required = required
			attr.Computed = computed
			attr.Optional = optional
			s.Attributes[e] = attr

		case schemaR.ListAttribute:
			attr := s.Attributes[e].(schemaR.ListAttribute)
			attr.Required = required
			attr.Computed = computed
			attr.Optional = optional
			s.Attributes[e] = attr

		case schemaR.SetAttribute:
			attr := s.Attributes[e].(schemaR.SetAttribute)
			attr.Required = required
			attr.Computed = computed
			attr.Optional = optional
			s.Attributes[e] = attr
		}
	}

	return s
}
