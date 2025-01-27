package org

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/org"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/org"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type CertificateLibraryModel struct {
	ID          supertypes.StringValue `tfsdk:"id"`
	Name        supertypes.StringValue `tfsdk:"name"`
	Description supertypes.StringValue `tfsdk:"description"`
	Certificate supertypes.StringValue `tfsdk:"certificate"`
	PrivateKey  supertypes.StringValue `tfsdk:"private_key"`
	Passphrase  supertypes.StringValue `tfsdk:"passphrase"`
}

type CertificateLibraryDatasourceModel struct {
	ID          supertypes.StringValue `tfsdk:"id"`
	Name        supertypes.StringValue `tfsdk:"name"`
	Description supertypes.StringValue `tfsdk:"description"`
	Certificate supertypes.StringValue `tfsdk:"certificate"`
}

// ToSDKCertificateLibraryModel converts the data Terraform to the SDK model.
func (rm *CertificateLibraryModel) ToSDKCertificateLibraryModel() org.CertificateModel {
	return org.CertificateModel{
		ID:          rm.ID.Get(),
		Name:        rm.Name.Get(),
		Description: rm.Description.Get(),
		Certificate: rm.Certificate.Get(),
		PrivateKey:  rm.PrivateKey.Get(),
		Passphrase:  rm.Passphrase.Get(),
	}
}

// Copy returns a copy of the CertificateLibraryModel.
func (rm *CertificateLibraryModel) Copy() *CertificateLibraryModel {
	x := &CertificateLibraryModel{}
	utils.ModelCopy(rm, x)
	return x
}
