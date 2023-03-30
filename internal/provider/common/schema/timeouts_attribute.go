package superschema

import (
	"context"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

var _ Attribute = TimeoutAttribute{}

type ResourceTimeoutAttribute struct {
	Create bool
	Read   bool
	Delete bool
	Update bool
}

type DatasourceTimeoutAttribute struct {
	Read bool
}
type TimeoutAttribute struct {
	Resource   *ResourceTimeoutAttribute
	DataSource *DatasourceTimeoutAttribute
}

// IsResource returns true if the attribute is a resource attribute.
func (s TimeoutAttribute) IsResource() bool {
	return s.Resource != nil
}

// IsDataSource returns true if the attribute is a data source attribute.
func (s TimeoutAttribute) IsDataSource() bool {
	return s.DataSource != nil
}

func (s TimeoutAttribute) GetResource(ctx context.Context) schemaR.Attribute {
	var a schemaR.Attribute

	if s.Resource != nil {
		a = timeoutsR.Attributes(ctx, timeoutsR.Opts{
			Create: s.Resource.Create,
			Read:   s.Resource.Read,
			Delete: s.Resource.Delete,
			Update: s.Resource.Update,
		})
	}
	return a
}

func (s TimeoutAttribute) GetDataSource(ctx context.Context) schemaD.Attribute {
	var a schemaD.Attribute

	if s.DataSource != nil && s.DataSource.Read {
		a = timeoutsD.Attributes(ctx)
	}
	return a
}
