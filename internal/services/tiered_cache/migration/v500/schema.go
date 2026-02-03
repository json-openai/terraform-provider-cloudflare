package v500

import (
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// SourceTieredCacheSchema returns the source schema for legacy cloudflare_tiered_cache resource.
// Schema version: 0 (SDKv2 default - v4 provider didn't set explicit version)
// Resource type: cloudflare_tiered_cache
//
// This minimal schema is used only for reading v4 state during migration.
// It includes only the properties needed for state parsing (Required, Optional, Computed, ElementType).
// Validators, Descriptions are intentionally kept minimal.
func SourceTieredCacheSchema() schema.Schema {
	return schema.Schema{
		Version: 0, // v4 provider schema version (SDKv2 default)
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"zone_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown(), stringplanmodifier.RequiresReplace()},
			},
			"cache_type": schema.StringAttribute{
				Required: true,
			},
			"editable": schema.BoolAttribute{
				Computed: true,
			},
			"modified_on": schema.StringAttribute{
				Computed:   true,
				CustomType: timetypes.RFC3339Type{},
			},
		},
	}
}
