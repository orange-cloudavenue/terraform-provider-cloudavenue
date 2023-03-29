package superschema

import (
	"context"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type Attributes map[string]Attribute

func (a Attributes) process(ctx context.Context, s schemaType) any {
	switch s {
	case resource:
		attributes := make(map[string]schemaR.Attribute)

		for k, v := range a {
			if v.IsResource() {
				attributes[k] = v.GetResource(ctx)
			}
		}
		return attributes

	case dataSource:
		attributes := make(map[string]schemaD.Attribute)

		for k, v := range a {
			if v.IsDataSource() {
				attributes[k] = v.GetDataSource(ctx)
			}
		}
		return attributes
	}

	return nil
}

type Attribute interface {
	IsResource() bool
	IsDataSource() bool
	GetResource(ctx context.Context) schemaR.Attribute
	GetDataSource(ctx context.Context) schemaD.Attribute
}
