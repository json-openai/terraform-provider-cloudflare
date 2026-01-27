package tiered_cache_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/acctest"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/consts"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// TestMigrateTieredCache_Smart tests migration from v4 cache_type = "smart" to v5 value = "on"
// NOTE: The tf-migrate tool creates BOTH cloudflare_tiered_cache AND cloudflare_argo_tiered_caching
// resources. The argo_tiered_caching resource will be created on next apply (this is by design).
func TestMigrateTieredCache_Smart(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	tieredCacheResourceName := "cloudflare_tiered_cache." + rnd
	argoTieredCachingResourceName := "cloudflare_argo_tiered_caching." + rnd
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	tmpDir := t.TempDir()

	// V4 config using cache_type = "smart"
	v4Config := fmt.Sprintf(`
resource "cloudflare_tiered_cache" "%[1]s" {
  zone_id    = "%[2]s"
  cache_type = "smart"
}`, rnd, zoneID)

	stateChecks := []statecheck.StateCheck{
		// Verify tiered_cache state transformation
		statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
		statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New("value"), knownvalue.StringExact("on")),
		statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New("id"), knownvalue.StringExact(zoneID)),
		// Verify new argo_tiered_caching resource was created
		statecheck.ExpectKnownValue(argoTieredCachingResourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
		statecheck.ExpectKnownValue(argoTieredCachingResourceName, tfjsonpath.New("value"), knownvalue.StringExact("on")),
		statecheck.ExpectKnownValue(argoTieredCachingResourceName, tfjsonpath.New("id"), knownvalue.StringExact(zoneID)),
	}

	// Build steps: Step 1 creates with v4 provider
	steps := []resource.TestStep{
		{
			ExternalProviders: map[string]resource.ExternalProvider{
				"cloudflare": {
					Source:            "cloudflare/cloudflare",
					VersionConstraint: "4.52.1",
				},
			},
			Config: v4Config,
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
				statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New("cache_type"), knownvalue.StringExact("smart")),
				statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New("id"), knownvalue.StringExact(zoneID)),
			},
		},
	}

	// Steps 2-3: Run migration and verify state (allows creates for split resources)
	steps = append(steps, acctest.MigrationV2TestStepAllowCreate(t, v4Config, tmpDir, "4.52.1", "v4", "v5", stateChecks)...)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
		},
		CheckDestroy: nil,
		WorkingDir:   tmpDir,
		Steps:        steps,
	})
}

// TestMigrateTieredCache_Generic tests migration from v4 cache_type = "generic"
// NOTE: The tf-migrate tool creates BOTH cloudflare_tiered_cache (with value="off") AND
// cloudflare_argo_tiered_caching (with value="on") resources.
func TestMigrateTieredCache_Generic(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	tieredCacheResourceName := "cloudflare_tiered_cache." + rnd
	argoTieredCachingResourceName := "cloudflare_argo_tiered_caching." + rnd
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	tmpDir := t.TempDir()

	// V4 config using cache_type = "generic"
	v4Config := fmt.Sprintf(`
resource "cloudflare_tiered_cache" "%[1]s" {
  zone_id    = "%[2]s"
  cache_type = "generic"
}`, rnd, zoneID)

	stateChecks := []statecheck.StateCheck{
		// Verify tiered_cache state transformation (value="off" for generic)
		statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
		statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New("value"), knownvalue.StringExact("off")),
		statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New("id"), knownvalue.StringExact(zoneID)),
		// Verify new argo_tiered_caching resource was created (value="on")
		statecheck.ExpectKnownValue(argoTieredCachingResourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
		statecheck.ExpectKnownValue(argoTieredCachingResourceName, tfjsonpath.New("value"), knownvalue.StringExact("on")),
		statecheck.ExpectKnownValue(argoTieredCachingResourceName, tfjsonpath.New("id"), knownvalue.StringExact(zoneID)),
	}

	// Build steps: Step 1 creates with v4 provider
	steps := []resource.TestStep{
		{
			ExternalProviders: map[string]resource.ExternalProvider{
				"cloudflare": {
					Source:            "cloudflare/cloudflare",
					VersionConstraint: "4.52.1",
				},
			},
			Config: v4Config,
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
				statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New("cache_type"), knownvalue.StringExact("generic")),
				statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New("id"), knownvalue.StringExact(zoneID)),
			},
		},
	}

	// Steps 2-3: Run migration and verify state (allows creates for split resources)
	steps = append(steps, acctest.MigrationV2TestStepAllowCreate(t, v4Config, tmpDir, "4.52.1", "v4", "v5", stateChecks)...)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
		},
		CheckDestroy: nil,
		WorkingDir:   tmpDir,
		Steps:        steps,
	})
}

// TestMigrateTieredCache_Off tests migration from v4 cache_type = "off" to v5 value = "off"
// NOTE: The tf-migrate tool creates BOTH cloudflare_tiered_cache AND cloudflare_argo_tiered_caching
// resources with value="off". The argo_tiered_caching resource will be created on next apply.
func TestMigrateTieredCache_Off(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	tieredCacheResourceName := "cloudflare_tiered_cache." + rnd
	argoTieredCachingResourceName := "cloudflare_argo_tiered_caching." + rnd
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	tmpDir := t.TempDir()

	// V4 config using cache_type = "off"
	v4Config := fmt.Sprintf(`
resource "cloudflare_tiered_cache" "%[1]s" {
  zone_id    = "%[2]s"
  cache_type = "off"
}`, rnd, zoneID)

	stateChecks := []statecheck.StateCheck{
		// Verify tiered_cache state transformation
		statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
		statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New("value"), knownvalue.StringExact("off")),
		statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New("id"), knownvalue.StringExact(zoneID)),
		// Verify new argo_tiered_caching resource was created (value="off")
		statecheck.ExpectKnownValue(argoTieredCachingResourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
		statecheck.ExpectKnownValue(argoTieredCachingResourceName, tfjsonpath.New("value"), knownvalue.StringExact("off")),
		statecheck.ExpectKnownValue(argoTieredCachingResourceName, tfjsonpath.New("id"), knownvalue.StringExact(zoneID)),
	}

	// Build steps: Step 1 creates with v4 provider
	steps := []resource.TestStep{
		{
			ExternalProviders: map[string]resource.ExternalProvider{
				"cloudflare": {
					Source:            "cloudflare/cloudflare",
					VersionConstraint: "4.52.1",
				},
			},
			Config: v4Config,
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
				statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New("cache_type"), knownvalue.StringExact("off")),
				statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New("id"), knownvalue.StringExact(zoneID)),
			},
		},
	}

	// Steps 2-3: Run migration and verify state (allows creates for split resources)
	steps = append(steps, acctest.MigrationV2TestStepAllowCreate(t, v4Config, tmpDir, "4.52.1", "v4", "v5", stateChecks)...)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
		},
		CheckDestroy: nil,
		WorkingDir:   tmpDir,
		Steps:        steps,
	})
}

// TestMigrateTieredCache_AllValues tests value transformations for all cache_type values
// NOTE: The tf-migrate tool creates BOTH cloudflare_tiered_cache AND cloudflare_argo_tiered_caching
// resources for each original tiered_cache resource.
func TestMigrateTieredCache_AllValues(t *testing.T) {
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")

	testCases := []struct {
		name               string
		v4Value            string
		v5TieredCacheValue string
		v5ArgoTieredValue  string
		description        string
	}{
		{
			name:               "Smart_To_On",
			v4Value:            "smart",
			v5TieredCacheValue: "on",
			v5ArgoTieredValue:  "on",
			description:        "Migration from smart: both resources get value=on",
		},
		{
			name:               "Off_To_Off",
			v4Value:            "off",
			v5TieredCacheValue: "off",
			v5ArgoTieredValue:  "off",
			description:        "Migration from off: both resources get value=off",
		},
		{
			name:               "Generic_Split",
			v4Value:            "generic",
			v5TieredCacheValue: "off",
			v5ArgoTieredValue:  "on",
			description:        "Migration from generic: tiered_cache=off, argo_tiered_caching=on",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rnd := utils.GenerateRandomResourceName()
			tieredCacheResourceName := "cloudflare_tiered_cache." + rnd
			argoTieredCachingResourceName := "cloudflare_argo_tiered_caching." + rnd
			tmpDir := t.TempDir()

			v4Config := fmt.Sprintf(`
resource "cloudflare_tiered_cache" "%[1]s" {
  zone_id    = "%[2]s"
  cache_type = "%[3]s"
}`, rnd, zoneID, tc.v4Value)

			stateChecks := []statecheck.StateCheck{
				// Verify tiered_cache state transformation
				statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
				statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New("value"), knownvalue.StringExact(tc.v5TieredCacheValue)),
				statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New("id"), knownvalue.StringExact(zoneID)),
				// Verify new argo_tiered_caching resource was created
				statecheck.ExpectKnownValue(argoTieredCachingResourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
				statecheck.ExpectKnownValue(argoTieredCachingResourceName, tfjsonpath.New("value"), knownvalue.StringExact(tc.v5ArgoTieredValue)),
				statecheck.ExpectKnownValue(argoTieredCachingResourceName, tfjsonpath.New("id"), knownvalue.StringExact(zoneID)),
			}

			// Build steps: Step 1 creates with v4 provider
			steps := []resource.TestStep{
				{
					ExternalProviders: map[string]resource.ExternalProvider{
						"cloudflare": {
							Source:            "cloudflare/cloudflare",
							VersionConstraint: "4.52.1",
						},
					},
					Config: v4Config,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
						statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New("cache_type"), knownvalue.StringExact(tc.v4Value)),
						statecheck.ExpectKnownValue(tieredCacheResourceName, tfjsonpath.New("id"), knownvalue.StringExact(zoneID)),
					},
				},
			}

			// Steps 2-3: Run migration and verify state (allows creates for split resources)
			steps = append(steps, acctest.MigrationV2TestStepAllowCreate(t, v4Config, tmpDir, "4.52.1", "v4", "v5", stateChecks)...)

			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					acctest.TestAccPreCheck(t)
				},
				WorkingDir: tmpDir,
				Steps:      steps,
			})
		})
	}
}

// TestMigrateTieredCache_V5VersionUpgrade tests upgrading between v5 versions (no-op upgrade)
// This ensures that existing v5 state is compatible with the latest provider version
func TestMigrateTieredCache_V5VersionUpgrade(t *testing.T) {
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")

	// Test different v5 versions to ensure state compatibility
	testCases := []struct {
		name    string
		version string
	}{
		{
			name:    "from_v5_8_0",
			version: "5.8.0",
		},
		{
			name:    "from_v5_10_0",
			version: "5.10.0",
		},
		{
			name:    "from_v5_12_0",
			version: "5.12.0",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rnd := utils.GenerateRandomResourceName()
			resourceName := "cloudflare_tiered_cache." + rnd
			tmpDir := t.TempDir()

			// V5 config using value attribute
			v5Config := fmt.Sprintf(`
resource "cloudflare_tiered_cache" "%[1]s" {
  zone_id = "%[2]s"
  value   = "on"
}`, rnd, zoneID)

			resource.Test(t, resource.TestCase{
				PreCheck: func() {
					acctest.TestAccPreCheck(t)
				},
				WorkingDir: tmpDir,
				Steps: []resource.TestStep{
					{
						// Step 1: Create with specific v5 provider version
						ExternalProviders: map[string]resource.ExternalProvider{
							"cloudflare": {
								Source:            "cloudflare/cloudflare",
								VersionConstraint: tc.version,
							},
						},
						Config: v5Config,
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue(resourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
							statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("value"), knownvalue.StringExact("on")),
							statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.StringExact(zoneID)),
						},
					},
					{
						// Step 2: Upgrade to latest provider - should be a no-op
						ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
						Config:                   v5Config,
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownValue(resourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
							statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("value"), knownvalue.StringExact("on")),
							statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.StringExact(zoneID)),
							statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("editable"), knownvalue.Bool(true)),
						},
					},
				},
			})
		})
	}
}

// TestMigrateTieredCache_V5VersionUpgrade_Off tests v5 version upgrade with value="off"
func TestMigrateTieredCache_V5VersionUpgrade_Off(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	resourceName := "cloudflare_tiered_cache." + rnd
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	tmpDir := t.TempDir()

	// V5 config with value = "off"
	v5Config := fmt.Sprintf(`
resource "cloudflare_tiered_cache" "%[1]s" {
  zone_id = "%[2]s"
  value   = "off"
}`, rnd, zoneID)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
		},
		WorkingDir: tmpDir,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with v5.8.0 provider
				ExternalProviders: map[string]resource.ExternalProvider{
					"cloudflare": {
						Source:            "cloudflare/cloudflare",
						VersionConstraint: "5.8.0",
					},
				},
				Config: v5Config,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("value"), knownvalue.StringExact("off")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.StringExact(zoneID)),
				},
			},
			{
				// Step 2: Upgrade to latest provider - should be a no-op
				ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
				Config:                   v5Config,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("value"), knownvalue.StringExact("off")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.StringExact(zoneID)),
				},
			},
		},
	})
}

// TestMigrateTieredCache_V5VersionUpgrade_WithVariables tests v5 version upgrade with variable references
func TestMigrateTieredCache_V5VersionUpgrade_WithVariables(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	resourceName := "cloudflare_tiered_cache." + rnd
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	tmpDir := t.TempDir()

	// V5 config with variable reference
	v5Config := fmt.Sprintf(`
variable "zone_id" {
  type    = string
  default = "%[2]s"
}

resource "cloudflare_tiered_cache" "%[1]s" {
  zone_id = var.zone_id
  value   = "on"
}`, rnd, zoneID)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
		},
		WorkingDir: tmpDir,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with v5.8.0 provider
				ExternalProviders: map[string]resource.ExternalProvider{
					"cloudflare": {
						Source:            "cloudflare/cloudflare",
						VersionConstraint: "5.8.0",
					},
				},
				Config: v5Config,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("value"), knownvalue.StringExact("on")),
				},
			},
			{
				// Step 2: Upgrade to latest provider - should be a no-op
				ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
				Config:                   v5Config,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New(consts.ZoneIDSchemaKey), knownvalue.StringExact(zoneID)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("value"), knownvalue.StringExact("on")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("editable"), knownvalue.Bool(true)),
				},
			},
		},
	})
}
