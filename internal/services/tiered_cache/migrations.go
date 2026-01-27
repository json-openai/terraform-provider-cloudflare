// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package tiered_cache

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithUpgradeState = (*TieredCacheResource)(nil)

func (r *TieredCacheResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   priorSchemaV0(),
			StateUpgrader: upgradeTieredCacheStateV0toV1,
		},
	}
}

// priorSchemaV0 returns the schema for version 0 (v4 provider schema with cache_type)
func priorSchemaV0() *schema.Schema {
	return &schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:   "Identifier.",
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"zone_id": schema.StringAttribute{
				Description:   "Identifier.",
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"cache_type": schema.StringAttribute{
				Description: "The type of tiered cache to use. Available values: smart, generic, off.",
				Required:    true,
			},
			"editable": schema.BoolAttribute{
				Description: "Whether the setting is editable.",
				Computed:    true,
			},
			"modified_on": schema.StringAttribute{
				Description: "Last time this setting was modified.",
				Computed:    true,
				CustomType:  timetypes.RFC3339Type{},
			},
		},
	}
}

// tieredCacheModelV0 represents the v4 state structure
type tieredCacheModelV0 struct {
	ID         types.String      `tfsdk:"id"`
	ZoneID     types.String      `tfsdk:"zone_id"`
	CacheType  types.String      `tfsdk:"cache_type"`
	Editable   types.Bool        `tfsdk:"editable"`
	ModifiedOn timetypes.RFC3339 `tfsdk:"modified_on"`
}

// upgradeTieredCacheStateV0toV1 upgrades the state from version 0 (v4) to version 1 (v5)
// This converts cache_type to value:
// - "smart" -> "on"
// - "off" -> "off"
// - "generic" -> error (should be moved to cloudflare_argo_tiered_caching via MoveState)
func upgradeTieredCacheStateV0toV1(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	// Get the old state
	var oldState tieredCacheModelV0
	diags := req.State.Get(ctx, &oldState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Transform cache_type to value
	var newValue string
	cacheType := oldState.CacheType.ValueString()
	switch cacheType {
	case "smart":
		newValue = "on"
	case "off":
		newValue = "off"
	case "generic":
		// generic should be handled via MoveState to cloudflare_argo_tiered_caching
		resp.Diagnostics.AddError(
			"Invalid State Upgrade",
			fmt.Sprintf("Cannot upgrade cloudflare_tiered_cache with cache_type='generic'. "+
				"Resources with cache_type='generic' should be migrated to cloudflare_argo_tiered_caching "+
				"using the tf-migrate tool. Please run the migration tool first."),
		)
		return
	default:
		// For unknown values (e.g., variable references), preserve as-is
		// The value will be validated by the schema validator
		newValue = cacheType
	}

	// Create the upgraded state
	upgradedState := TieredCacheModel{
		ID:         oldState.ID,
		ZoneID:     oldState.ZoneID,
		Value:      types.StringValue(newValue),
		Editable:   oldState.Editable,
		ModifiedOn: oldState.ModifiedOn,
	}

	// Set the upgraded state
	diags = resp.State.Set(ctx, upgradedState)
	resp.Diagnostics.Append(diags...)
}
