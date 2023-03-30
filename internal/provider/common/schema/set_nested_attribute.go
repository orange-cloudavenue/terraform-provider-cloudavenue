package superschema

import (
	"context"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var _ Attribute = SetNestedAttribute{}

type SetNestedAttribute struct {
	Common     *schemaR.SetNestedAttribute
	Resource   *schemaR.SetNestedAttribute
	DataSource *schemaD.SetNestedAttribute
	Attributes Attributes
}

// IsResource returns true if the attribute is a resource attribute.
func (s SetNestedAttribute) IsResource() bool {
	return s.Resource != nil || s.Common != nil
}

// IsDataSource returns true if the attribute is a data source attribute.
func (s SetNestedAttribute) IsDataSource() bool {
	return s.DataSource != nil || s.Common != nil
}

func (s SetNestedAttribute) GetResource(ctx context.Context) schemaR.Attribute {
	a := schemaR.SetNestedAttribute{
		NestedObject: schemaR.NestedAttributeObject{
			Attributes: s.Attributes.process(ctx, resource).(map[string]schemaR.Attribute),
		},
	}

	if s.Common != nil {
		a.Required = s.Common.Required
		a.Optional = s.Common.Optional
		a.Computed = s.Common.Computed
		a.Sensitive = s.Common.Sensitive
		a.MarkdownDescription = s.Common.MarkdownDescription
		a.Description = s.Common.Description
		a.DeprecationMessage = s.Common.DeprecationMessage
		a.Validators = s.Common.Validators
		a.PlanModifiers = s.Common.PlanModifiers
		a.Default = s.Common.Default
	}

	//nolint:dupl
	if s.Resource != nil {
		if s.Resource.Required {
			a.Required = true
		}

		if s.Resource.Optional {
			a.Optional = true
		}

		if s.Resource.Computed {
			a.Computed = true
		}

		if s.Resource.Sensitive {
			a.Sensitive = true
		}

		if s.Resource.MarkdownDescription != "" {
			a.MarkdownDescription += s.Resource.MarkdownDescription
		}

		if s.Resource.Description != "" {
			a.Description += s.Resource.Description
		}

		if s.Resource.DeprecationMessage != "" {
			a.DeprecationMessage += s.Resource.DeprecationMessage
		}

		if len(s.Resource.Validators) > 0 {
			a.Validators = append(a.Validators, s.Resource.Validators...)
		}

		if len(s.Resource.PlanModifiers) > 0 {
			a.PlanModifiers = append(a.PlanModifiers, s.Resource.PlanModifiers...)
		}

		if s.Resource.Default != nil {
			a.Default = s.Resource.Default
		}
	}

	return a
}

func (s SetNestedAttribute) GetDataSource(ctx context.Context) schemaD.Attribute {
	a := schemaD.SetNestedAttribute{
		NestedObject: schemaD.NestedAttributeObject{
			Attributes: s.Attributes.process(ctx, dataSource).(map[string]schemaD.Attribute),
		},
	}

	if s.Common != nil {
		a.Required = s.Common.Required
		a.Optional = s.Common.Optional
		a.Computed = s.Common.Computed
		a.Sensitive = s.Common.Sensitive
		a.MarkdownDescription = s.Common.MarkdownDescription
		a.Description = s.Common.Description
		a.DeprecationMessage = s.Common.DeprecationMessage
		a.Validators = s.Common.Validators
	}

	if s.DataSource != nil {
		if s.DataSource.Required {
			a.Required = true
		}

		if s.DataSource.Optional {
			a.Optional = true
		}

		if s.DataSource.Computed {
			a.Computed = true
		}

		if s.DataSource.Sensitive {
			a.Sensitive = true
		}

		if s.DataSource.MarkdownDescription != "" {
			a.MarkdownDescription += s.DataSource.MarkdownDescription
		}

		if s.DataSource.Description != "" {
			a.Description += s.DataSource.Description
		}

		if s.DataSource.DeprecationMessage != "" {
			a.DeprecationMessage += s.DataSource.DeprecationMessage
		}

		if len(s.DataSource.Validators) > 0 {
			a.Validators = append(a.Validators, s.DataSource.Validators...)
		}
	}

	return a
}
