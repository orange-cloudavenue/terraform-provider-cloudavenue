package storageprofile

import (
	"errors"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
)

var (
	ErrStorageProfileNameIsEmpty = errors.New("storage profile Name is empty")
	ErrStorageProfileNotFound    = errors.New("storage profile not found")
)

type Handler interface {
	GetStorageProfile(storageProfileName string, refresh bool) (*govcdtypes.CatalogStorageProfiles, error)
	GetStorageProfileReference(storageProfileName string, refresh bool) (*govcdtypes.Reference, error)
	FindStorageProfileName(storageProfileName string) (string, error)
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

func SuperSchema() superschema.StringAttribute {
	return superschema.StringAttribute{
		Common: &schemaR.StringAttribute{
			MarkdownDescription: "The storage profile name to use.",
			Computed:            true,
		},
		Resource: &schemaR.StringAttribute{
			Optional: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf(storageProfileValues...),
			},
		},
	}
}
