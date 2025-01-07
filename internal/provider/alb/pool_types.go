/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package alb

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type albPoolModel struct {
	ID                       types.String `tfsdk:"id"`
	EdgeGatewayID            types.String `tfsdk:"edge_gateway_id"`
	EdgeGatewayName          types.String `tfsdk:"edge_gateway_name"`
	Name                     types.String `tfsdk:"name"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	Description              types.String `tfsdk:"description"`
	Algorithm                types.String `tfsdk:"algorithm"`
	DefaultPort              types.Int64  `tfsdk:"default_port"`
	GracefulTimeoutPeriod    types.Int64  `tfsdk:"graceful_timeout_period"`
	Members                  types.Set    `tfsdk:"members"`
	HealthMonitors           types.Set    `tfsdk:"health_monitors"`
	PersistenceProfile       types.Object `tfsdk:"persistence_profile"`
	PassiveMonitoringEnabled types.Bool   `tfsdk:"passive_monitoring_enabled"`

	// CACertificateIDs         types.Set    `tfsdk:"ca_certificate_ids"`
	// CNCheckEnabled           types.Bool   `tfsdk:"cn_check_enabled"`
	// DomainNames              types.Set    `tfsdk:"domain_names"`
}

type member struct {
	Enabled   types.Bool   `tfsdk:"enabled"`
	IPAddress types.String `tfsdk:"ip_address"`
	Port      types.Int64  `tfsdk:"port"`
	Ratio     types.Int64  `tfsdk:"ratio"`
}

var memberAttrTypes = map[string]attr.Type{
	"enabled":    types.BoolType,
	"ip_address": types.StringType,
	"port":       types.Int64Type,
	"ratio":      types.Int64Type,
}

type persistenceProfile struct {
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

var persistenceProfileAttrTypes = map[string]attr.Type{
	"type":  types.StringType,
	"value": types.StringType,
}
