/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegw

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func vpnIPSecSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "Provides a resource to manage an IPsec VPN Tunnel. You can configure a site-to-site connectivity between an Edge Gateway and remote site. The remote site must support IPSec protocol. The VPN is able to initiate and respond to incoming tunnel requests. The VPN tunnel is established only when both sides of the tunnel are configured. The VPN tunnel is terminated when one side of the tunnel is deleted or disabled. The VPN tunnel is re-established when the disabled side is enabled again.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "Provides a data source to read IPsec VPN Tunnel configuration of your Edge Gateway.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the IPsec VPN Tunnel Configuration.",
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Name of the IPsec VPN Tunnel Configuration.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Optional: true,
					Computed: true,
				},
			},
			"edge_gateway_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Edge Gateway",
					Computed:            true,
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"edge_gateway_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Name of the Edge Gateway.",
					Computed:            true,
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"description": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A description of the IPsec VPN Tunnel Configuration.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"enabled": superschema.SuperBoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Enable or Disable the IPsec VPN Tunnel Configuration.",
					Computed:            true,
				},
				Resource: &schemaR.BoolAttribute{
					Default:  booldefault.StaticBool(true),
					Optional: true,
				},
			},
			"pre_shared_key": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Pre-Shared Key (PSK) is an Authentication method. Is a complex password (ASCII) that will be exchanged between both sites in order to set up the IPsec tunnel.",
					Sensitive:           true,
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"local_ip_address": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "An IPv4 Address for the local endpoint. This has to be a sub-allocated IP on the Edge Gateway. This endpoint must be reach by the remote endpoint.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						fstringvalidator.IsIP(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"local_networks": superschema.SuperSetAttribute{
				Common: &schemaR.SetAttribute{
					MarkdownDescription: "Set of local networks in CIDR format. This local_networks will be exchanged between both sites in order to route ip traffic in VPN tunnel.",
					ElementType:         supertypes.StringType{},
				},
				Resource: &schemaR.SetAttribute{
					Required: true,
				},
				DataSource: &schemaD.SetAttribute{
					Computed: true,
				},
			},
			"remote_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Remote ID is an optional identity used to establish the VPN tunnel. If not set, the Remote IP Address will be used as Remote ID.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
			},
			"remote_ip_address": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "An IPv4 Address for the remote endpoint. This is your remote VPN endpoint you need to reach.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						fstringvalidator.IsIP(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"remote_networks": superschema.SuperSetAttribute{
				Common: &schemaR.SetAttribute{
					MarkdownDescription: "Set of remote networks in CIDR format. This remote_networks will be exchanged between both sites in order to route ip traffic in VPN tunnel.",
					ElementType:         supertypes.StringType{},
				},
				Resource: &schemaR.SetAttribute{
					Required: true,
				},
				DataSource: &schemaD.SetAttribute{
					Computed: true,
				},
			},
			"security_type": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Type of Security Profile used for the IPsec VPN Tunnel.",
					Computed:            true,
				},
			},
			"security_profile": superschema.SuperSingleNestedAttribute{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Customization of your IPSec configuration. The configuration used must be symmetric for both endpoint VPN.",
					Computed:            true,
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
				},
				Attributes: map[string]superschema.Attribute{
					"ike_version": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "IKE (Internet Key Exchange) is an encrypt protocol of your VPN data.",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								fstringvalidator.OneOfWithDescription(
									fstringvalidator.OneOfWithDescriptionValues{
										Value:       "IKE_V1",
										Description: "When you select this option, IPSec VPN initiates and responds to IKEv1 protocol only.",
									},
									fstringvalidator.OneOfWithDescriptionValues{
										Value:       "IKE_V2",
										Description: "The default option. When you select this version, IPSec VPN initiates and responds to IKEv2 protocol only.",
									},
									fstringvalidator.OneOfWithDescriptionValues{
										Value:       "IKE_FLEX",
										Description: "When you select this option, if the tunnel establishment fails with IKEv2 protocol, the source site does not fall back and initiate a connection with the IKEv1 protocol. Instead, if the remote site initiates a connection with the IKEv1 protocol, then the connection is accepted.",
									},
								),
							},
						},
					},
					"ike_encryption_algorithm": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Encryption algorithms used by IKE.",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("AES_128", "AES_256", "AES_GCM_128", "AES_GCM_192", "AES_GCM_256"),
							},
						},
					},
					"ike_digest_algorithm": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Secure hashing algorithms to use during the IKE negotiation.",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("SHA1", "SHA2_256", "SHA2_384", "SHA2_512"),
								fstringvalidator.NullIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("ike_encryption_algorithm"), []attr.Value{types.StringValue("AES_GCM_128"), types.StringValue("AES_GCM_256"), types.StringValue("AES_GCM_512")}),
							},
						},
					},
					"ike_dh_groups": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The Diffie-Hellman (DH) key exchange algorithm is a method used to make a shared encryption key available to two entities over an insecure communications channel.",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("GROUP2", "GROUP5", "GROUP14", "GROUP15", "GROUP16", "GROUP19", "GROUP20", "GROUP21"),
							},
						},
					},
					"ike_sa_lifetime": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "Security association lifetime in seconds. It is number of seconds before the IPsec tunnel ike part needs to reestablish.",
							Computed:            true,
						},
						Resource: &schemaR.Int64Attribute{
							Optional: true,
							Default:  int64default.StaticInt64(86400),
							Validators: []validator.Int64{
								int64validator.Between(21600, 31536000),
							},
						},
					},
					"tunnel_pfs": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "PFS (Perfect Forward Secrecy) capacity enabled or disabled. It's generates unique private keys for each secure session.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional:   true,
							Default:    booldefault.StaticBool(true),
							Validators: []validator.Bool{
								// TODO - Issue open https://github.com/orange-cloudavenue/terraform-plugin-framework-validators/issues/88
							},
						},
					},
					"tunnel_df_policy": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Policy for handling defragmentation.",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Default:  stringdefault.StaticString("COPY"),
							Validators: []validator.String{
								fstringvalidator.OneOfWithDescription(
									fstringvalidator.OneOfWithDescriptionValues{
										Value:       "COPY",
										Description: "Copies the defragmentation bit from the inner IP packet to the outer packet.",
									},
									fstringvalidator.OneOfWithDescriptionValues{
										Value:       "CLEAR",
										Description: "Ignores the defragmentation bit present in the inner packet.",
									},
								),
							},
						},
					},
					"tunnel_encryption_algorithms": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Encryption algorithms to use in IPSec tunnel establishment.",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("AES_128", "AES_256", "AES_GCM_128", "AES_GCM_192", "AES_GCM_256"),
							},
						},
					},
					"tunnel_digest_algorithms": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Digest algorithms to be used for message digest.",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("SHA1", "SHA2_256", "SHA2_384", "SHA2_512"),
								fstringvalidator.NullIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("tunnel_encryption_algorithms"), []attr.Value{types.StringValue("AES_GCM_128"), types.StringValue("AES_GCM_256"), types.StringValue("AES_GCM_512")}),
							},
						},
					},
					"tunnel_dh_groups": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The Diffie-Hellman (DH) key exchange algorithm is a method used to make a shared encryption key available to two entities over an insecure communications channel.",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("GROUP2", "GROUP5", "GROUP14", "GROUP15", "GROUP16", "GROUP19", "GROUP20", "GROUP21"),
							},
						},
					},
					"tunnel_sa_lifetime": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "Security association lifetime in seconds. It is number of seconds before the IPsec tunnel needs to reestablish.",
							Computed:            true,
						},
						Resource: &schemaR.Int64Attribute{
							Optional: true,
							Default:  int64default.StaticInt64(3600),
							Validators: []validator.Int64{
								int64validator.Between(900, 31536000),
							},
						},
					},
					"tunnel_dpd": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "Value in seconds of Dead Probe Detection interval.",
							Computed:            true,
						},
						Resource: &schemaR.Int64Attribute{
							Optional: true,
							Default:  int64default.StaticInt64(60),
							Validators: []validator.Int64{
								int64validator.Between(3, 60),
							},
						},
					},
				},
			},
		},
	}
}
