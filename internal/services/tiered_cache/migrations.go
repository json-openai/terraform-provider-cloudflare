package tiered_cache

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/services/tiered_cache/migration/v500"
)

var _ resource.ResourceWithUpgradeState = (*TieredCacheResource)(nil)

// UpgradeState registers state upgraders for schema version changes.
//
// This handles two upgrade paths:
// 1. v4 state (schema_version=0) → v5 (version=500): Full transformation (cache_type → value)
// 2. v5 state (version=1) → v5 (version=500): No-op upgrade (when TF_MIG_TEST=1)
//
// The separation of schema versions (v4=0, v5=1/500) eliminates the need for
// dual-format detection that was required in earlier implementations.
func (r *TieredCacheResource) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	sourceSchema := v500.SourceTieredCacheSchema()
	targetSchema := ResourceSchema(ctx)

	return map[int64]resource.StateUpgrader{
		// Handle state from v4 SDKv2 provider (schema_version=0)
		// Transforms cache_type to value: "smart"→"on", "off"→"off"
		// "generic" returns error (must use tf-migrate tool to split)
		0: {
			PriorSchema:   &sourceSchema,
			StateUpgrader: v500.UpgradeFromV0,
		},

		// Handle state from v5 Plugin Framework provider with version=1
		// This is a no-op upgrade that just bumps the version to 500
		1: {
			PriorSchema:   &targetSchema,
			StateUpgrader: v500.UpgradeFromV1,
		},
	}
}
