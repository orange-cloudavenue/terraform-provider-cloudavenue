package diskparams

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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

const descriptionStorageProfile = "Storage profile to override the VM default one."

// StorageProfileAttribute returns the schema.Attribute for the storage profile.
func StorageProfileAttribute() schema.Attribute {
	return schema.StringAttribute{
		MarkdownDescription: descriptionStorageProfile + " Allowed values are: " + storageProfileValuesDescription + ".",
		Computed:            true,
		Optional:            true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		Validators: []validator.String{
			stringvalidator.OneOf(storageProfileValues...),
		},
	}
}

// StorageProfileAttributeComputed returns the schema.Attribute for the storage profile.
func StorageProfileAttributeComputed() schema.Attribute {
	return schema.StringAttribute{
		MarkdownDescription: descriptionStorageProfile,
		Computed:            true,
	}
}

// StorageProfileAttributeRequired returns the schema.Attribute for the storage profile.
func StorageProfileAttributeRequired() schema.Attribute {
	return schema.StringAttribute{
		MarkdownDescription: descriptionStorageProfile + " Allowed values are: " + storageProfileValuesDescription + ".",
		Required:            true,
		Validators: []validator.String{
			stringvalidator.OneOf(storageProfileValues...),
		},
	}
}
