package catalog

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
)

const (
	categoryName = "catalog"

	schemaID   = "catalog_id"
	schemaName = "catalog_name"
)

type catalog interface {
	GetID() string
	GetName() string
	GetIDOrName() string
	GetCatalog() (*govcd.AdminCatalog, error)
}

type base struct {
	id   string
	name string
}

/*
schemaCatalogName

	returns the schema.Attribute for the catalog name.

	Default values are :
	- Optional: false
	- Computed: false
	- Required: true

	You can override the default values by using the following options:
	- IsComputed()
	- IsRequired()
	- IsOptional()

	If the override is define all the default values are set to false.
*/
func schemaCatalogName(opts ...common.AttributeOpts) schema.Attribute {
	// Initialize the attribute options.
	a := &common.AttributeStruct{}

	// if opts is empty, set the default values.
	if len(opts) == 0 {
		a.Required = true
	} else {
		// Override the default values with the provided options.
		for _, opt := range opts {
			opt(a)
		}
	}

	description := "The name of the catalog."
	if a.Optional {
		description += " Required if `catalog_id` is not set."
	}

	sAttribute := schema.StringAttribute{
		MarkdownDescription: description,
		Computed:            a.Computed,
		Optional:            a.Optional,
		Required:            a.Required,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}

	if a.Optional {
		sAttribute.Validators = []validator.String{
			stringvalidator.ExactlyOneOf(path.MatchRoot("catalog_id"), path.MatchRoot("catalog_name")),
		}
	}

	return sAttribute
}

/*
schemaCatalogID

	returns the schema.Attribute for the catalog id.

	Default values are :
	- Optional: false
	- Computed: false
	- Required: true

	You can override the default values by using the following options:
	- IsComputed()
	- IsRequired()
	- IsOptional()

	If the override is define all the default values are set to false.
	If the attribute is optional, the validator is set to check if one of the two attributes is set.
*/
func schemaCatalogID(opts ...common.AttributeOpts) schema.Attribute {
	// Initialize the attribute options.
	a := &common.AttributeStruct{}

	// if opts is empty, set the default values.
	if len(opts) == 0 {
		a.Required = true
	} else {
		// Override the default values with the provided options.
		for _, opt := range opts {
			opt(a)
		}
	}

	description := "The ID of the catalog."
	if a.Optional {
		description += " Required if `catalog_name` is not set."
	}

	sAttribute := schema.StringAttribute{
		MarkdownDescription: description,
		Computed:            a.Computed,
		Optional:            a.Optional,
		Required:            a.Required,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}

	if a.Optional {
		sAttribute.Validators = []validator.String{
			stringvalidator.ExactlyOneOf(path.MatchRoot("catalog_id"), path.MatchRoot("catalog_name")),
		}
	}

	return sAttribute
}
