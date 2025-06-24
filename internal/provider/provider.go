/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package provider provides the CloudAvenue Terraform Provider.
package provider

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	casdk "github.com/orange-cloudavenue/cloudavenue-sdk-go"
	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	clientnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/netbackup"
	caverrors "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider          = &cloudavenueProvider{}
	_ tfprotov6.ProviderServer   = &cloudavenueProvider{}
	_ tfprotov6.ResourceServer   = &cloudavenueProvider{}
	_ tfprotov6.DataSourceServer = &cloudavenueProvider{}
)

// cloudavenueProvider is the provider implementation.
type cloudavenueProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &cloudavenueProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *cloudavenueProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cloudavenue"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *cloudavenueProvider) Schema(ctx context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerSchema(ctx)
}

func (p *cloudavenueProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config cloudavenueProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudAvenue := client.CloudAvenue{
		// This is a new SDK Cloudavenue
		CAVSDKOpts: &casdk.ClientOpts{
			Netbackup: &clientnetbackup.Opts{
				Endpoint: emptyOrValue(config.NetBackupURL),
				Username: emptyOrValue(config.NetBackupUser),
				Password: emptyOrValue(config.NetBackupPassword),
			},
			CloudAvenue: &clientcloudavenue.Opts{
				URL:      emptyOrValue(config.URL),
				Username: emptyOrValue(config.User),
				Password: emptyOrValue(config.Password),
				Org:      emptyOrValue(config.Org),
				VDC:      emptyOrValue(config.VDC),
			},
		},
	}

	const (
		summaryErrorAPICAV = "Unable to Create Cloud Avenue API Client"
		summaryErrorVCD    = "Unable to Create VMWare VCD Client"
	)

	cA, err := cloudAvenue.New()
	if err != nil {
		switch {
		case errors.Is(err, caverrors.ErrEmpty):
			resp.Diagnostics.AddError("Attribute is empty", err.Error())
			return
		case errors.Is(err, caverrors.ErrInvalidFormat):
			resp.Diagnostics.AddError("Attribute has invalid format", err.Error())
			return
		default:
			resp.Diagnostics.AddError(summaryErrorAPICAV, "unknown error: "+err.Error())
			return
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Make the CloudAvenue client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = cA
	resp.ResourceData = cA
}

func emptyOrValue(value basetypes.StringValue) string {
	if value.IsNull() {
		return ""
	}
	return value.ValueString()
}

// ! tfprotov6.ProviderServer interface methods
//
// GetResourceIdentitySchemas returns the schemas for the provider's resource identity types.
func (p *cloudavenueProvider) GetResourceIdentitySchemas(_ context.Context, req *tfprotov6.GetResourceIdentitySchemasRequest) (*tfprotov6.GetResourceIdentitySchemasResponse, error) {
	resp := &tfprotov6.GetResourceIdentitySchemasResponse{}
	return resp, nil
}

// UpgradeResourceIdentity upgrades the resource identity type for the provider's resources.
func (p *cloudavenueProvider) UpgradeResourceIdentity(_ context.Context, req *tfprotov6.UpgradeResourceIdentityRequest) (*tfprotov6.UpgradeResourceIdentityResponse, error) {
	resp := &tfprotov6.UpgradeResourceIdentityResponse{}
	return resp, nil
}

// ApplyResourceChange applies the changes to the provider's resources.
func (p *cloudavenueProvider) ApplyResourceChange(_ context.Context, req *tfprotov6.ApplyResourceChangeRequest) (*tfprotov6.ApplyResourceChangeResponse, error) {
	resp := &tfprotov6.ApplyResourceChangeResponse{}
	return resp, nil
}

// ImportResourceState imports the state of the provider's resources.
func (p *cloudavenueProvider) ImportResourceState(_ context.Context, req *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
	resp := &tfprotov6.ImportResourceStateResponse{}
	return resp, nil
}

// CallFunction is a no-op for the provider.
func (p *cloudavenueProvider) CallFunction(_ context.Context, req *tfprotov6.CallFunctionRequest) (*tfprotov6.CallFunctionResponse, error) {
	resp := &tfprotov6.CallFunctionResponse{}
	return resp, nil
}

// MoveResourceState moves the state of the provider's resources.
func (p *cloudavenueProvider) MoveResourceState(_ context.Context, req *tfprotov6.MoveResourceStateRequest) (*tfprotov6.MoveResourceStateResponse, error) {
	resp := &tfprotov6.MoveResourceStateResponse{}
	return resp, nil
}

// CloseEphemeralResource closes the ephemeral resource for the provider.
func (p *cloudavenueProvider) CloseEphemeralResource(_ context.Context, req *tfprotov6.CloseEphemeralResourceRequest) (*tfprotov6.CloseEphemeralResourceResponse, error) {
	resp := &tfprotov6.CloseEphemeralResourceResponse{}
	return resp, nil
}

// OpenEphemeralResource opens the ephemeral resource for the provider.
func (p *cloudavenueProvider) OpenEphemeralResource(_ context.Context, req *tfprotov6.OpenEphemeralResourceRequest) (*tfprotov6.OpenEphemeralResourceResponse, error) {
	resp := &tfprotov6.OpenEphemeralResourceResponse{}
	return resp, nil
}

// ConfigureProvider configures the provider with the given configuration.
func (p *cloudavenueProvider) ConfigureProvider(ctx context.Context, req *tfprotov6.ConfigureProviderRequest) (*tfprotov6.ConfigureProviderResponse, error) {
	resp := &tfprotov6.ConfigureProviderResponse{}
	return resp, nil
}

// GetFunction gets the function for the provider.
func (p *cloudavenueProvider) GetFunctions(_ context.Context, req *tfprotov6.GetFunctionsRequest) (*tfprotov6.GetFunctionsResponse, error) {
	resp := &tfprotov6.GetFunctionsResponse{}
	return resp, nil
}

// PlanResourceChange plans the changes to the provider's resources.
func (p *cloudavenueProvider) PlanResourceChange(_ context.Context, req *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	resp := &tfprotov6.PlanResourceChangeResponse{}
	return resp, nil
}

// ReadResource reads the state of the provider's resources.
func (p *cloudavenueProvider) ReadResource(_ context.Context, req *tfprotov6.ReadResourceRequest) (*tfprotov6.ReadResourceResponse, error) {
	resp := &tfprotov6.ReadResourceResponse{}
	return resp, nil
}

// ReadDataSource reads the state of the provider's data sources.
func (p *cloudavenueProvider) ReadDataSource(_ context.Context, req *tfprotov6.ReadDataSourceRequest) (*tfprotov6.ReadDataSourceResponse, error) {
	resp := &tfprotov6.ReadDataSourceResponse{}
	return resp, nil
}

// UpgradeResourceState upgrades the state of the provider's resources.
func (p *cloudavenueProvider) UpgradeResourceState(_ context.Context, req *tfprotov6.UpgradeResourceStateRequest) (*tfprotov6.UpgradeResourceStateResponse, error) {
	resp := &tfprotov6.UpgradeResourceStateResponse{}
	return resp, nil
}

// ValidateResourceConfig validates the configuration of the provider's resources.
func (p *cloudavenueProvider) ValidateResourceConfig(_ context.Context, req *tfprotov6.ValidateResourceConfigRequest) (*tfprotov6.ValidateResourceConfigResponse, error) {
	resp := &tfprotov6.ValidateResourceConfigResponse{}
	return resp, nil
}

// ValidateDataResourceConfig validates the configuration of the provider's data sources.
func (p *cloudavenueProvider) ValidateDataResourceConfig(_ context.Context, req *tfprotov6.ValidateDataResourceConfigRequest) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	resp := &tfprotov6.ValidateDataResourceConfigResponse{}
	return resp, nil
}

// GetMetadata returns the metadata for the provider.
func (p *cloudavenueProvider) GetMetadata(_ context.Context, req *tfprotov6.GetMetadataRequest) (*tfprotov6.GetMetadataResponse, error) {
	resp := &tfprotov6.GetMetadataResponse{}
	return resp, nil
}

// GetProviderSchema returns the provider schema.
func (p *cloudavenueProvider) GetProviderSchema(_ context.Context, req *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	resp := &tfprotov6.GetProviderSchemaResponse{}
	return resp, nil
}

// RenewEphemeralResource renews the ephemeral resource for the provider.
func (p *cloudavenueProvider) RenewEphemeralResource(_ context.Context, req *tfprotov6.RenewEphemeralResourceRequest) (*tfprotov6.RenewEphemeralResourceResponse, error) {
	resp := &tfprotov6.RenewEphemeralResourceResponse{}
	return resp, nil
}

// StopProvider stops the provider.
func (p *cloudavenueProvider) StopProvider(_ context.Context, req *tfprotov6.StopProviderRequest) (*tfprotov6.StopProviderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}

// ValidateEphemeralResourceConfig validates the configuration of the ephemeral resource for the provider.
func (p *cloudavenueProvider) ValidateEphemeralResourceConfig(_ context.Context, req *tfprotov6.ValidateEphemeralResourceConfigRequest) (*tfprotov6.ValidateEphemeralResourceConfigResponse, error) {
	resp := &tfprotov6.ValidateEphemeralResourceConfigResponse{}
	return resp, nil
}

// ValidateProviderConfig validates the configuration of the provider.
func (p *cloudavenueProvider) ValidateProviderConfig(_ context.Context, req *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	resp := &tfprotov6.ValidateProviderConfigResponse{}
	return resp, nil
}
