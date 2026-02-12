package v500

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// UpgradeNoOp handles state upgrades within the v5 series (schema_version=1+).
// This is a no-op upgrade since the schema is compatible - just copy state through.
func UpgradeNoOp(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	tflog.Info(ctx, "Upgrading pages_project state (no-op)")
	// No-op upgrade: schema is compatible, just copy raw state through
	resp.State.Raw = req.State.Raw
}

// UpgradeFromV0 upgrades the state from v4 (version 0) to v5.
// v4 SDKv2 stores blocks as lists, so we access the first element of each slice.
func UpgradeFromV0(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	var priorState SourcePagesProjectModelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Transform the v0 (v4 provider) state to v5 state
	newState, diags := Transform(ctx, priorState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}
