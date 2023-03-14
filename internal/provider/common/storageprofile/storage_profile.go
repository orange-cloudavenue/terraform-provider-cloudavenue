package storageprofile

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common"
)

var (
	ErrStorageProfileIDIsEmpty = errors.New("storage profile ID is empty")
	ErrStorageProfileNotFound  = errors.New("storage profile not found")
)

type Handler interface {
	GetStorageProfile(storageProfileID string, refresh bool) (*govcdtypes.CatalogStorageProfiles, error)
	GetStorageProfileReference(storageProfileID string, refresh bool) (*govcdtypes.Reference, error)
	FindStorageProfileID(storageProfileName string) (string, error)
}

const (
	SchemaStorageProfile = "storage_profile"
)

type storageProfile string

// String.
func (s storageProfile) String() string {
	return string(s)
}

const (
	storageProfileSilver      storageProfile = "silver"
	storageProfileSilverR1    storageProfile = "silver_r1"
	storageProfileSilverR2    storageProfile = "silver_r2"
	storageProfileGold        storageProfile = "gold"
	storageProfileGoldR1      storageProfile = "gold_r1"
	storageProfileGoldR2      storageProfile = "gold_r2"
	storageProfileGoldHM      storageProfile = "gold_hm"
	storageProfilePlatinum3   storageProfile = "platinum3k"
	storageProfilePlatinum3R1 storageProfile = "platinum3k_r1"
	storageProfilePlatinum3R2 storageProfile = "platinum3k_r2"
	storageProfilePlatinum3HM storageProfile = "platinum3k_hm"
	storageProfilePlatinum7   storageProfile = "platinum7k"
	storageProfilePlatinum7R1 storageProfile = "platinum7k_r1"
	storageProfilePlatinum7R2 storageProfile = "platinum7k_r2"
	storageProfilePlatinum7HM storageProfile = "platinum7k_hm"
)

var storageProfileValues = []string{
	storageProfileSilver.String(),
	storageProfileSilverR1.String(),
	storageProfileSilverR2.String(),
	storageProfileGold.String(),
	storageProfileGoldR1.String(),
	storageProfileGoldR2.String(),
	storageProfileGoldHM.String(),
	storageProfilePlatinum3.String(),
	storageProfilePlatinum3R1.String(),
	storageProfilePlatinum3R2.String(),
	storageProfilePlatinum3HM.String(),
	storageProfilePlatinum7.String(),
	storageProfilePlatinum7R1.String(),
	storageProfilePlatinum7R2.String(),
	storageProfilePlatinum7HM.String(),
}

var storageProfileValuesDescription = func() string {
	var s string
	countItems := len(storageProfileValues)
	for i, v := range storageProfileValues {
		if i == countItems-1 {
			s += "`" + v + "`"
		} else {
			s += "`" + v + "`, "
		}
	}
	return s
}()

/*
Schema

	returns the schema.Attribute for the storage profile.

	Default values are :
	- Optional: true
	- Computed: false
	- Required: false

	You can override the default values by using the following options:
	- IsComputed()
	- IsRequired()
	- IsOptional()

	If the override is define all the default values are set to false.
*/
func Schema(opts ...common.AttributeOpts) schema.Attribute {
	// Initialize the attribute options.
	a := &common.AttributeStruct{}

	// if opts is empty, set the default values.
	if len(opts) == 0 {
		a.Optional = true
	} else {
		// Override the default values with the provided options.
		for _, opt := range opts {
			opt(a)
		}
	}

	description := "Storage profile to override the VM default one."
	if a.Optional || a.Required {
		description += " Allowed values are: " + storageProfileValuesDescription + "."
	}

	return schema.StringAttribute{
		MarkdownDescription: description,
		Computed:            a.Computed,
		Optional:            a.Optional,
		Required:            a.Required,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf(storageProfileValues...),
		},
	}
}
