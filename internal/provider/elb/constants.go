/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package elb

// Shared attribute name constants for ELB schemas.
const (
	edgeGatewayID   = "edge_gateway_id"
	edgeGatewayName = "edge_gateway_name"
	enabled         = "enabled"
)

// Shared attribute name constants for ELB HTTP policies schemas.
const (
	virtualServiceID = "virtual_service_id"
	policies         = "policies"
	name             = "name"
	active           = "active"
	logging          = "logging"
	criteria         = "criteria"
	protocol         = "protocol"
	clientIP         = "client_ip"
	ipAddresses      = "ip_addresses"
	servicePorts     = "service_ports"
	ports            = "ports"
	httpMethods      = "http_methods"
	methods          = "methods"
	pathAttr         = "path"
	paths            = "paths"
	cookie           = "cookie"
	value            = "value"
	requestHeaders   = "request_headers"
	values           = "values"
	query            = "query"
	actions          = "actions"
	host             = "host"
	keepQuery        = "keep_query"
	port             = "port"
	statusCode       = "status_code"
)

// Additional shared attribute name constants for security/request/response policies schemas.
const (
	redirect        = "redirect"
	modifyHeaders   = "modify_headers"
	rewriteURL      = "rewrite_url"
	connection      = "connection"
	rateLimit       = "rate_limit"
	sendResponse    = "send_response"
	localResponse   = "local_response"
	closeConnection = "close_connection"
	redirectToHTTPS = "redirect_to_https"
)

// Shared MarkdownDescription constants for ELB HTTP policies schemas.
const (
	clientIPDescription       = "Match the rule based on client IP address rules."
	criteriaDescription       = "Criteria to match."
	ipAddressesDescription    = "IP addresses to match."
	servicePortsDescription   = "Match the rule based on service port rules."
	portsDescription          = "A port list allows you to define which service ports (e.g.: [80, 443] ) the HTTP security policy should match."
	httpMethodsDescription    = "Match the rule based on HTTP method rules."
	methodsDescription        = "Methods to match."
	pathDescription           = "Match the rule based on path rules."
	pathsDescription          = "A set of paths to match given criteria."
	cookieDescription         = "Match the rule based on cookie rules."
	cookieNameDescription     = "Name of the cookie to match."
	cookieValueDescription    = "Value of the cookie to match."
	requestHeadersDescription = "Match the rule based on request headers rules."
	headerNameDescription     = "Name of the HTTP header whose value is to be matched."
	headerValuesDescription   = "Values of the HTTP header to match."
	queryDescription          = "Text contained in the query string"
	actionsDescription        = "Actions to perform when the rule matches."
	hostDescription           = "Host to which redirect the request. Default is the original host"
	keepQueryDescription      = "Keep or drop the query of the incoming request URI in the redirected URI"
	redirectPathDescription   = "Path to which redirect the request. Default is the original path"
	portDescription           = "Port to which redirect the request."
	httpProtocolDescription   = "HTTP protocol"
)
