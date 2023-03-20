package superschema

import (
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

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
