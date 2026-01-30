// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package v500

import (
	"context"
	"strings"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Transform converts source (legacy SDKv2 v4) state to target (current Plugin Framework v5) state.
// This function handles the migration from cloudflare_load_balancer_pool in SDKv2 to Plugin Framework.
func Transform(ctx context.Context, source SourceCloudflareLoadBalancerPoolModel) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Validate required fields
	if source.AccountID.IsNull() || source.AccountID.IsUnknown() {
		diags.AddError("Missing required field", "account_id is required for load_balancer_pool migration")
		return nil, diags
	}
	if source.Name.IsNull() || source.Name.IsUnknown() {
		diags.AddError("Missing required field", "name is required for load_balancer_pool migration")
		return nil, diags
	}

	tflog.Debug(ctx, "Starting transformation", map[string]interface{}{
		"account_id": source.AccountID.ValueString(),
		"name":       source.Name.ValueString(),
	})

	target := &TargetLoadBalancerPoolModel{
		ID:                source.ID,
		AccountID:         source.AccountID,
		Name:              source.Name,
		Enabled:           source.Enabled,
		MinimumOrigins:    source.MinimumOrigins,
		Latitude:          source.Latitude,
		Longitude:         source.Longitude,
		Description:       source.Description,
		Monitor:           source.Monitor,
		NotificationEmail: source.NotificationEmail,
		CreatedOn:         source.CreatedOn,
		ModifiedOn:        source.ModifiedOn,
	}

	// New v5 fields - set to null/empty
	target.MonitorGroup = types.StringNull()
	target.DisabledAt = timetypes.NewRFC3339Null()
	target.Networks = customfield.NullList[types.String](ctx)
	target.NotificationFilter = customfield.NullObject[TargetLoadBalancerPoolNotificationFilterModel](ctx)

	// Transform check_regions: Set → List
	if !source.CheckRegions.IsNull() && !source.CheckRegions.IsUnknown() {
		checkRegionsList, checkDiags := convertSetToList(ctx, source.CheckRegions)
		diags.Append(checkDiags...)
		if !diags.HasError() {
			target.CheckRegions = &checkRegionsList
		}
	}

	// Transform origins: []SourceOriginsModel → *[]*TargetLoadBalancerPoolOriginsModel
	if len(source.Origins) > 0 {
		originsList, originsDiags := transformOrigins(ctx, source.Origins)
		diags.Append(originsDiags...)
		if !diags.HasError() {
			target.Origins = &originsList
		}
	}

	// Transform load_shedding: Array[0] → NestedObject
	// SDK v2 stores TypeList MaxItems:1 as array: [{...}]
	if len(source.LoadShedding) > 0 {
		loadSheddingObj, loadSheddingDiags := transformLoadShedding(ctx, source.LoadShedding[0])
		diags.Append(loadSheddingDiags...)
		if !diags.HasError() {
			target.LoadShedding = loadSheddingObj
		}
	} else {
		target.LoadShedding = customfield.NullObject[TargetLoadBalancerPoolLoadSheddingModel](ctx)
	}

	// Transform origin_steering: Array[0] → NestedObject
	// SDK v2 stores TypeList MaxItems:1 as array: [{...}]
	if len(source.OriginSteering) > 0 {
		originSteeringObj, originSteeringDiags := transformOriginSteering(ctx, source.OriginSteering[0])
		diags.Append(originSteeringDiags...)
		if !diags.HasError() {
			target.OriginSteering = originSteeringObj
		}
	} else {
		target.OriginSteering = customfield.NullObject[TargetLoadBalancerPoolOriginSteeringModel](ctx)
	}

	tflog.Debug(ctx, "Transformation completed", map[string]interface{}{
		"target_id":   target.ID.ValueString(),
		"origins_len": func() int { if target.Origins != nil { return len(*target.Origins) } else { return 0 } }(),
	})

	return target, diags
}

// convertSetToList converts a types.Set to a slice of types.String for List fields.
func convertSetToList(ctx context.Context, set types.Set) ([]types.String, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract to native Go []string first
	var rawStrings []string
	diags.Append(set.ElementsAs(ctx, &rawStrings, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Convert to []types.String
	result := make([]types.String, len(rawStrings))
	for i, str := range rawStrings {
		result[i] = types.StringValue(str)
	}

	return result, diags
}

// transformOrigins converts source origins (SDK v2 Set stored as array) to target origins (List of pointers).
func transformOrigins(ctx context.Context, sourceOrigins []SourceOriginsModel) ([]*TargetLoadBalancerPoolOriginsModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	targetOrigins := make([]*TargetLoadBalancerPoolOriginsModel, len(sourceOrigins))

	for i, sourceOrigin := range sourceOrigins {
		targetOrigin := &TargetLoadBalancerPoolOriginsModel{
			Address:          sourceOrigin.Address,
			Enabled:          sourceOrigin.Enabled,
			Name:             sourceOrigin.Name,
			VirtualNetworkID: sourceOrigin.VirtualNetworkID,
			Weight:           sourceOrigin.Weight,
		}

		// New v5 fields
		targetOrigin.Port = types.Int64Null()
		targetOrigin.DisabledAt = timetypes.NewRFC3339Null()

		// Transform header: Complex nested structure change
		if len(sourceOrigin.Header) > 0 {
			headerObj, headerDiags := transformOriginHeader(ctx, sourceOrigin.Header[0])
			diags.Append(headerDiags...)
			if !diags.HasError() {
				targetOrigin.Header = headerObj
			}
		}

		targetOrigins[i] = targetOrigin
	}

	return targetOrigins, diags
}

// transformOriginHeader transforms the header structure from v4 to v5 format.
//
// v4 format (SDK v2 TypeSet stored as array):
//
//	[{ header: "Host", values: Set["value1", "value2"] }]
//
// v5 format:
//
//	{ host: ["value1", "value2"] }
//
// Per user decision: Assume all v4 headers are "Host" (safest for migration).
func transformOriginHeader(ctx context.Context, sourceHeader SourceHeaderModel) (*TargetLoadBalancerPoolOriginsHeaderModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Convert header name to lowercase for v5 key
	// Per user decision: assume all headers are "Host", but handle gracefully
	headerName := sourceHeader.Header.ValueString()
	headerKey := strings.ToLower(headerName)

	tflog.Debug(ctx, "Transforming origin header", map[string]interface{}{
		"header_name": headerName,
		"header_key":  headerKey,
	})

	// Only support "host" in v5
	if headerKey != "host" {
		tflog.Warn(ctx, "Non-Host header found in v4 state, v5 only supports Host header", map[string]interface{}{
			"header_name": headerName,
		})
		// Still transform it, but as "host" since that's all v5 supports
	}

	// Convert values Set → List
	if sourceHeader.Values.IsNull() || sourceHeader.Values.IsUnknown() {
		return nil, diags
	}

	valuesList, valuesDiags := convertSetToList(ctx, sourceHeader.Values)
	diags.Append(valuesDiags...)
	if diags.HasError() {
		return nil, diags
	}

	targetHeader := &TargetLoadBalancerPoolOriginsHeaderModel{
		Host: &valuesList,
	}

	return targetHeader, diags
}

// transformLoadShedding converts source load_shedding (array element) to target (NestedObject).
//
// SDK v2 stores TypeList MaxItems:1 as: [{ default_percent: 0, default_policy: "", ... }]
// v5 stores as object: { default_percent: 0, default_policy: "random", ... }
//
// Per user decision: If v4 has empty string for policies, set v5 defaults.
func transformLoadShedding(ctx context.Context, source SourceLoadSheddingModel) (customfield.NestedObject[TargetLoadBalancerPoolLoadSheddingModel], diag.Diagnostics) {
	var diags diag.Diagnostics

	target := TargetLoadBalancerPoolLoadSheddingModel{
		DefaultPercent: source.DefaultPercent,
		SessionPercent: source.SessionPercent,
	}

	// Handle default_policy: v4 default "" → v5 default "random"
	if source.DefaultPolicy.IsNull() || source.DefaultPolicy.ValueString() == "" {
		target.DefaultPolicy = types.StringValue("random")
	} else {
		target.DefaultPolicy = source.DefaultPolicy
	}

	// Handle session_policy: v4 default "" → v5 default "hash"
	if source.SessionPolicy.IsNull() || source.SessionPolicy.ValueString() == "" {
		target.SessionPolicy = types.StringValue("hash")
	} else {
		target.SessionPolicy = source.SessionPolicy
	}

	tflog.Debug(ctx, "Transformed load_shedding", map[string]interface{}{
		"default_policy": target.DefaultPolicy.ValueString(),
		"session_policy": target.SessionPolicy.ValueString(),
	})

	nestedObj, objDiags := customfield.NewObject(ctx, &target)
	diags.Append(objDiags...)
	return nestedObj, diags
}

// transformOriginSteering converts source origin_steering (array element) to target (NestedObject).
//
// SDK v2 stores TypeList MaxItems:1 as: [{ policy: "random" }]
// v5 stores as object: { policy: "random" }
func transformOriginSteering(ctx context.Context, source SourceOriginSteeringModel) (customfield.NestedObject[TargetLoadBalancerPoolOriginSteeringModel], diag.Diagnostics) {
	var diags diag.Diagnostics

	target := TargetLoadBalancerPoolOriginSteeringModel{
		Policy: source.Policy,
	}

	// Default is "random" in both v4 and v5, so no special handling needed
	if target.Policy.IsNull() || target.Policy.ValueString() == "" {
		target.Policy = types.StringValue("random")
	}

	nestedObj, objDiags := customfield.NewObject(ctx, &target)
	diags.Append(objDiags...)
	return nestedObj, diags
}
