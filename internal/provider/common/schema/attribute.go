package superschema

import (
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type Attributes map[string]Attribute

func (a Attributes) process(s schemaType) any {
	switch s {
	case resource:
		attributes := make(map[string]schemaR.Attribute)

		for k, v := range a {
			if v.IsResource() {
				attributes[k] = v.GetResource()
			}
		}
		return attributes

	case dataSource:
		attributes := make(map[string]schemaD.Attribute)

		for k, v := range a {
			if v.IsDataSource() {
				attributes[k] = v.GetDataSource()
			}
		}
		return attributes
	}

	return nil
}

type Attribute interface {
	IsResource() bool
	IsDataSource() bool
	GetResource() schemaR.Attribute
	GetDataSource() schemaD.Attribute
}
