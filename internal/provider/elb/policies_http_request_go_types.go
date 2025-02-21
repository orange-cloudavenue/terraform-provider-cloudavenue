/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package elb

type (
	GoPoliciesHTTPRequestModel struct {
		ID               string                                `tfsdk:"id"`
		VirtualServiceID string                                `tfsdk:"virtual_service_id"`
		Policies         []*GoPoliciesHTTPRequestModelPolicies `tfsdk:"policies"`
	}

	GoPoliciesHTTPRequestModelPolicies struct {
		Name     string                              `tfsdk:"name"`
		Active   bool                                `tfsdk:"active"`
		Logging  bool                                `tfsdk:"logging"`
		Criteria *GoPoliciesHTTPRequestMatchCriteria `tfsdk:"criteria"`
		Actions  *GoPoliciesHTTPRequestActions       `tfsdk:"actions"`
	}

	GoPoliciesHTTPRequestMatchCriteria struct {
		Protocol       string                          `tfsdk:"protocol"`
		ClientIP       *GoPoliciesHTTPClientIPMatch    `tfsdk:"client_ip"`
		ServicePorts   *GoPoliciesHTTPServicePortMatch `tfsdk:"service_ports"`
		HTTPMethods    *GoPoliciesHTTPMethodMatch      `tfsdk:"http_methods"`
		Path           *GoPoliciesHTTPPathMatch        `tfsdk:"path"`
		Cookie         *GoPoliciesHTTPCookieMatch      `tfsdk:"cookie"`
		RequestHeaders []*GoPoliciesHTTPHeaderMatch    `tfsdk:"request_headers"`
		Query          []string                        `tfsdk:"query"`
	}

	GoPoliciesHTTPRequestActions struct {
		Redirect      *GoPoliciesHTTPActionRedirect        `tfsdk:"redirect"`
		RewriteURL    *GoPoliciesHTTPActionURLRewrite      `tfsdk:"rewrite_url"`
		ModifyHeaders []*GoPoliciesHTTPActionHeaderRewrite `tfsdk:"modify_headers"`
	}

	// * Match.
	GoPoliciesHTTPClientIPMatch struct {
		Criteria    string   `tfsdk:"criteria"`
		IPAddresses []string `tfsdk:"ip_addresses"`
	}
	GoPoliciesHTTPServicePortMatch struct {
		Criteria string  `tfsdk:"criteria"`
		Ports    []int64 `tfsdk:"ports"`
	}
	GoPoliciesHTTPMethodMatch struct {
		Criteria string   `tfsdk:"criteria"`
		Methods  []string `tfsdk:"methods"`
	}
	GoPoliciesHTTPPathMatch struct {
		Criteria string   `tfsdk:"criteria"`
		Paths    []string `tfsdk:"paths"`
	}
	GoPoliciesHTTPHeaderMatch struct {
		Criteria string   `tfsdk:"criteria"`
		Name     string   `tfsdk:"name"`
		Values   []string `tfsdk:"values"`
	}
	GoPoliciesHTTPCookieMatch struct {
		Criteria string `tfsdk:"criteria"`
		Name     string `tfsdk:"name"`
		Value    string `tfsdk:"value"`
	}

	// * Action.
	GoPoliciesHTTPActionRedirect struct {
		Host       string `tfsdk:"host"`
		KeepQuery  bool   `tfsdk:"keep_query"`
		Path       string `tfsdk:"path"`
		Port       *int   `tfsdk:"port"`
		Protocol   string `tfsdk:"protocol"`
		StatusCode int    `tfsdk:"status_code"`
	}
	GoPoliciesHTTPActionURLRewrite struct {
		Host      string `tfsdk:"host"`
		Path      string `tfsdk:"path"`
		Query     string `tfsdk:"query"`
		KeepQuery bool   `tfsdk:"keep_query"`
	}
	GoPoliciesHTTPActionHeaderRewrite struct {
		Action string `tfsdk:"action"`
		Name   string `tfsdk:"name"`
		Value  string `tfsdk:"value"`
	}
)
