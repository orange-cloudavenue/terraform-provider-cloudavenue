/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package elb

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type (
	policiesHTTPPrivateModel struct {
		EdgeGatewayID string `json:"edgeGatewayId"`
	}
)

type getEdgeGateway func(ctx context.Context, virtualServiceID string) (edgeGatewayID string, diags diag.Diagnostics)

type privateState interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

const privateStateKey = "private_data"

// Get retrieves the private data from the private state.
// If the private data is not found, it retrieves the edge gateway ID using the provided function and sets it in the private state.
// It also unmarshals the private data into the policiesHTTPPrivateModel struct.
// If the unmarshalling fails, it retrieves the edge gateway ID and sets it in the private state.
func (p *policiesHTTPPrivateModel) Get(ctx context.Context, virtualServiceID string, private privateState, gE getEdgeGateway) diag.Diagnostics {
	privateEncoded, d := private.GetKey(ctx, privateStateKey)
	if d.HasError() {
		goto getEdgeID
	}

	if err := json.Unmarshal(privateEncoded, p); err != nil {
		// Get the virtual service to check if it exists and retrieve the edge gateway ID for locking
		goto getEdgeID
	}
	goto END

getEdgeID:
	p.EdgeGatewayID, d = gE(ctx, virtualServiceID)
	if d.HasError() {
		return d
	}

END:
	return p.set(ctx, private)
}

func (p *policiesHTTPPrivateModel) set(ctx context.Context, private privateState) (diags diag.Diagnostics) {
	privateEncoded, err := json.Marshal(p)
	if err != nil {
		diags.AddError("Failed to marshal private data", err.Error())
		return diags
	}

	return private.SetKey(ctx, privateStateKey, privateEncoded)
}
