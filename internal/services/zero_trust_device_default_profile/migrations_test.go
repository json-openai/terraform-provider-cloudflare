package zero_trust_device_default_profile_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/acctest"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/utils"
)

// TestMigrateZeroTrustDeviceDefaultProfile_Basic tests migration of a basic default profile 
// Tests: resource rename, field removal (default), basic field preservation
func TestMigrateZeroTrustDeviceDefaultProfile_Basic(t *testing.T) {
	// Zero Trust resources require API_KEY + EMAIL, not API_TOKEN
	originalToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if originalToken != "" {
		os.Unsetenv("CLOUDFLARE_API_TOKEN")
		defer os.Setenv("CLOUDFLARE_API_TOKEN", originalToken)
	}

	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	rnd := utils.GenerateRandomResourceName()
	tmpDir := t.TempDir()

	// v4 config with cloudflare_zero_trust_device_profiles + default=true
	// Include service_mode_v2 to avoid API 500 errors with minimal configs
	v4Config := fmt.Sprintf(`
resource "cloudflare_zero_trust_device_profiles" "%[1]s" {
  account_id  = "%[2]s"
  name        = "Test Profile"
  description = "Test device profile for migration"
  default     = true

  allow_mode_switch    = true
  auto_connect         = 0
  captive_portal       = 180
  service_mode_v2_mode = "proxy"
  service_mode_v2_port = 8080
}`, rnd, accountID)

	stateChecks := []statecheck.StateCheck{
		// Verify resource type changed
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),

		// Verify fields preserved
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("allow_mode_switch"), knownvalue.Bool(true)),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("auto_connect"), knownvalue.Float64Exact(0)),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("captive_portal"), knownvalue.Float64Exact(180)),
	}

	migrationSteps := acctest.MigrationV2TestStepWithStateNormalization(t, v4Config, tmpDir, "4.52.1", "v4", "v5", stateChecks)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
			acctest.TestAccPreCheck_AccountID(t)
		},
		WorkingDir: tmpDir,
		Steps: append([]resource.TestStep{
			{
				// Step 1: Create with v4 provider
				ExternalProviders: map[string]resource.ExternalProvider{
					"cloudflare": {
						Source:            "cloudflare/cloudflare",
						VersionConstraint: "4.52.1",
					},
				},
				Config: v4Config,
			},
		}, migrationSteps...),
	})
}

// TestMigrateZeroTrustDeviceDefaultProfile_WithServiceModeV2 tests migration with service_mode_v2 transformation
// Tests: flat fields (service_mode_v2_mode, service_mode_v2_port) → nested object
func TestMigrateZeroTrustDeviceDefaultProfile_WithServiceModeV2(t *testing.T) {
	// Zero Trust resources require API_KEY + EMAIL, not API_TOKEN
	originalToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if originalToken != "" {
		os.Unsetenv("CLOUDFLARE_API_TOKEN")
		defer os.Setenv("CLOUDFLARE_API_TOKEN", originalToken)
	}

	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	rnd := utils.GenerateRandomResourceName()
	tmpDir := t.TempDir()

	// v4 config with service_mode_v2 flat fields
	// Note: Port only applies when mode is "proxy"
	v4Config := fmt.Sprintf(`
resource "cloudflare_zero_trust_device_profiles" "%[1]s" {
  account_id  = "%[2]s"
  name        = "Test Profile with Service Mode"
  description = "Test device profile with service_mode_v2 fields"
  default     = true

  allow_mode_switch    = true
  auto_connect         = 15
  captive_portal       = 300
  service_mode_v2_mode = "proxy"
  service_mode_v2_port = 8080
}`, rnd, accountID)

	stateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),

		// Verify service_mode_v2 nested object created
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("service_mode_v2").AtMapKey("mode"), knownvalue.StringExact("proxy")),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("service_mode_v2").AtMapKey("port"), knownvalue.Float64Exact(8080)),

		// Verify other fields
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("allow_mode_switch"), knownvalue.Bool(true)),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("auto_connect"), knownvalue.Float64Exact(15)),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("captive_portal"), knownvalue.Float64Exact(300)),
	}

	migrationSteps := acctest.MigrationV2TestStepWithStateNormalization(t, v4Config, tmpDir, "4.52.1", "v4", "v5", stateChecks)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
			acctest.TestAccPreCheck_AccountID(t)
		},
		WorkingDir: tmpDir,
		Steps: append([]resource.TestStep{
			{
				// Step 1: Create with v4 provider
				ExternalProviders: map[string]resource.ExternalProvider{
					"cloudflare": {
						Source:            "cloudflare/cloudflare",
						VersionConstraint: "4.52.1",
					},
				},
				Config: v4Config,
			},
		}, migrationSteps...),
	})
}

// TestMigrateZeroTrustDeviceDefaultProfile_Maximal tests migration with all optional fields
// Tests: comprehensive field preservation, type conversions
func TestMigrateZeroTrustDeviceDefaultProfile_Maximal(t *testing.T) {
	// Zero Trust resources require API_KEY + EMAIL, not API_TOKEN
	originalToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if originalToken != "" {
		os.Unsetenv("CLOUDFLARE_API_TOKEN")
		defer os.Setenv("CLOUDFLARE_API_TOKEN", originalToken)
	}

	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	rnd := utils.GenerateRandomResourceName()
	tmpDir := t.TempDir()

	// v4 config with all optional fields
	// Note: Port only applies when mode is "proxy"
	v4Config := fmt.Sprintf(`
resource "cloudflare_zero_trust_device_profiles" "%[1]s" {
  account_id  = "%[2]s"
  name        = "Maximal Test Profile"
  description = "Test device profile with all optional fields"
  default     = true

  allow_mode_switch     = false
  allow_updates         = true
  allowed_to_leave      = true
  auto_connect          = 30
  captive_portal        = 600
  disable_auto_fallback = true
  switch_locked         = true
  support_url           = "https://support.cf-tf-test.com"

  service_mode_v2_mode = "proxy"
  service_mode_v2_port = 443
}`, rnd, accountID)

	stateChecks := []statecheck.StateCheck{
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),

		// Verify boolean fields
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("allow_mode_switch"), knownvalue.Bool(false)),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("allow_updates"), knownvalue.Bool(true)),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("allowed_to_leave"), knownvalue.Bool(true)),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("disable_auto_fallback"), knownvalue.Bool(true)),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("switch_locked"), knownvalue.Bool(true)),

		// Verify numeric fields (converted to Float64)
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("auto_connect"), knownvalue.Float64Exact(30)),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("captive_portal"), knownvalue.Float64Exact(600)),

		// Verify string fields
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("support_url"), knownvalue.StringExact("https://support.cf-tf-test.com")),

		// Verify service_mode_v2 nested object
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("service_mode_v2").AtMapKey("mode"), knownvalue.StringExact("proxy")),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("service_mode_v2").AtMapKey("port"), knownvalue.Float64Exact(443)),
	}

	migrationSteps := acctest.MigrationV2TestStepWithStateNormalization(t, v4Config, tmpDir, "4.52.1", "v4", "v5", stateChecks)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
			acctest.TestAccPreCheck_AccountID(t)
		},
		WorkingDir: tmpDir,
		Steps: append([]resource.TestStep{
			{
				// Step 1: Create with v4 provider
				ExternalProviders: map[string]resource.ExternalProvider{
					"cloudflare": {
						Source:            "cloudflare/cloudflare",
						VersionConstraint: "4.52.1",
					},
				},
				Config: v4Config,
			},
		}, migrationSteps...),
	})
}

// TestMigrateZeroTrustDeviceDefaultProfile_OldResourceName tests migration from old resource name
// Tests: cloudflare_device_settings_policy (old name) → cloudflare_zero_trust_device_default_profile
func TestMigrateZeroTrustDeviceDefaultProfile_OldResourceName(t *testing.T) {
	// Zero Trust resources require API_KEY + EMAIL, not API_TOKEN
	originalToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if originalToken != "" {
		os.Unsetenv("CLOUDFLARE_API_TOKEN")
		defer os.Setenv("CLOUDFLARE_API_TOKEN", originalToken)
	}

	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	rnd := utils.GenerateRandomResourceName()
	tmpDir := t.TempDir()

	// v4 config with OLD resource name (cloudflare_device_settings_policy)
	// Include service_mode_v2 to avoid API 500 errors with minimal configs
	v4Config := fmt.Sprintf(`
resource "cloudflare_device_settings_policy" "%[1]s" {
  account_id  = "%[2]s"
  name        = "Old Name Test Profile"
  description = "Test device profile with old resource name"
  default     = true

  allow_mode_switch    = true
  auto_connect         = 0
  captive_portal       = 300
  service_mode_v2_mode = "proxy"
  service_mode_v2_port = 8080
}`, rnd, accountID)

	stateChecks := []statecheck.StateCheck{
		// Verify resource type changed to new name
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("allow_mode_switch"), knownvalue.Bool(true)),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("auto_connect"), knownvalue.Float64Exact(0)),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("captive_portal"), knownvalue.Float64Exact(300)),
	}

	migrationSteps := acctest.MigrationV2TestStepWithStateNormalization(t, v4Config, tmpDir, "4.52.1", "v4", "v5", stateChecks)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
			acctest.TestAccPreCheck_AccountID(t)
		},
		WorkingDir: tmpDir,
		Steps: append([]resource.TestStep{
			{
				// Step 1: Create with v4 provider using old name
				ExternalProviders: map[string]resource.ExternalProvider{
					"cloudflare": {
						Source:            "cloudflare/cloudflare",
						VersionConstraint: "4.52.1",
					},
				},
				Config: v4Config,
			},
		}, migrationSteps...),
	})
}

// TestMigrateZeroTrustDeviceDefaultProfile_WithSplitTunnelExclude tests split tunnel merge (exclude mode)
// Tests: cloudflare_split_tunnel (no policy_id, mode=exclude) merged into default profile's exclude array
func TestMigrateZeroTrustDeviceDefaultProfile_WithSplitTunnelExclude(t *testing.T) {
	// Zero Trust resources require API_KEY + EMAIL, not API_TOKEN
	originalToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if originalToken != "" {
		os.Unsetenv("CLOUDFLARE_API_TOKEN")
		defer os.Setenv("CLOUDFLARE_API_TOKEN", originalToken)
	}

	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	rnd := utils.GenerateRandomResourceName()
	tmpDir := t.TempDir()

	// v4 config with default profile + split tunnel (no policy_id = applies to default)
	// Use same minimal config as TestMigrateZeroTrustDeviceDefaultProfile_Basic
	v4Config := fmt.Sprintf(`
resource "cloudflare_zero_trust_device_profiles" "%[1]s" {
  account_id  = "%[2]s"
  name        = "Test Profile"
  description = "Test device profile for migration"
  default     = true

  allow_mode_switch    = true
  auto_connect         = 0
  captive_portal       = 180
  service_mode_v2_mode = "proxy"
  service_mode_v2_port = 8080
}

resource "cloudflare_split_tunnel" "%[1]s_exclude" {
  account_id = "%[2]s"
  mode       = "exclude"

  tunnels {
    address     = "192.168.1.0/24"
    description = "Local network"
  }

  tunnels {
    host        = "internal.local"
    description = "Internal domain"
  }
}`, rnd, accountID)

	stateChecks := []statecheck.StateCheck{
		// Verify default profile exists
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),

		// Verify split tunnel was merged into exclude array (SetNestedAttribute so order may vary)
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("exclude"), knownvalue.SetSizeExact(2)),

		// Note: Can't check exact position due to Set semantics, but verify fields are present
		// The migration should have merged the tunnels into the exclude array
	}

	migrationSteps := acctest.MigrationV2TestStepWithStateNormalization(t, v4Config, tmpDir, "4.52.1", "v4", "v5", stateChecks)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
			acctest.TestAccPreCheck_AccountID(t)
		},
		WorkingDir: tmpDir,
		Steps: append([]resource.TestStep{
			{
				// Step 1: Create with v4 provider
				ExternalProviders: map[string]resource.ExternalProvider{
					"cloudflare": {
						Source:            "cloudflare/cloudflare",
						VersionConstraint: "4.52.1",
					},
				},
				Config: v4Config,
			},
		}, migrationSteps...),
	})
}

// TestMigrateZeroTrustDeviceDefaultProfile_WithSplitTunnelInclude tests split tunnel merge (include mode)
// Tests: cloudflare_split_tunnel (no policy_id, mode=include) merged into default profile's include array
func TestMigrateZeroTrustDeviceDefaultProfile_WithSplitTunnelInclude(t *testing.T) {
	// Zero Trust resources require API_KEY + EMAIL, not API_TOKEN
	originalToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if originalToken != "" {
		os.Unsetenv("CLOUDFLARE_API_TOKEN")
		defer os.Setenv("CLOUDFLARE_API_TOKEN", originalToken)
	}

	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	rnd := utils.GenerateRandomResourceName()
	tmpDir := t.TempDir()

	// v4 config with default profile + split tunnel include mode
	// Use same minimal config as TestMigrateZeroTrustDeviceDefaultProfile_Basic
	v4Config := fmt.Sprintf(`
resource "cloudflare_zero_trust_device_profiles" "%[1]s" {
  account_id  = "%[2]s"
  name        = "Test Profile"
  description = "Test device profile for migration"
  default     = true

  allow_mode_switch    = true
  auto_connect         = 0
  captive_portal       = 180
  service_mode_v2_mode = "proxy"
  service_mode_v2_port = 8080
}

resource "cloudflare_split_tunnel" "%[1]s_include" {
  account_id = "%[2]s"
  mode       = "include"

  tunnels {
    address     = "203.0.113.0/24"
    description = "Corporate VPN"
  }
}`, rnd, accountID)

	stateChecks := []statecheck.StateCheck{
		// Verify default profile exists
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),

		// Verify split tunnel was merged into include array
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("include"), knownvalue.SetSizeExact(1)),
	}

	migrationSteps := acctest.MigrationV2TestStepWithStateNormalization(t, v4Config, tmpDir, "4.52.1", "v4", "v5", stateChecks)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
			acctest.TestAccPreCheck_AccountID(t)
		},
		WorkingDir: tmpDir,
		Steps: append([]resource.TestStep{
			{
				// Step 1: Create with v4 provider
				ExternalProviders: map[string]resource.ExternalProvider{
					"cloudflare": {
						Source:            "cloudflare/cloudflare",
						VersionConstraint: "4.52.1",
					},
				},
				Config: v4Config,
			},
		}, migrationSteps...),
	})
}

// TestMigrateZeroTrustDeviceDefaultProfile_WithMultipleSplitTunnels tests multiple split tunnels merge
// Tests: Multiple cloudflare_split_tunnel resources (both include and exclude) merged into default profile
func TestMigrateZeroTrustDeviceDefaultProfile_WithMultipleSplitTunnels(t *testing.T) {
	// Zero Trust resources require API_KEY + EMAIL, not API_TOKEN
	originalToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if originalToken != "" {
		os.Unsetenv("CLOUDFLARE_API_TOKEN")
		defer os.Setenv("CLOUDFLARE_API_TOKEN", originalToken)
	}

	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	rnd := utils.GenerateRandomResourceName()
	tmpDir := t.TempDir()

	// v4 config with default profile + multiple split tunnels (both modes)
	// Use same minimal config as TestMigrateZeroTrustDeviceDefaultProfile_Basic
	v4Config := fmt.Sprintf(`
resource "cloudflare_zero_trust_device_profiles" "%[1]s" {
  account_id  = "%[2]s"
  name        = "Test Profile"
  description = "Test device profile for migration"
  default     = true

  allow_mode_switch    = true
  auto_connect         = 0
  captive_portal       = 180
  service_mode_v2_mode = "proxy"
  service_mode_v2_port = 8080
}

resource "cloudflare_split_tunnel" "%[1]s_exclude" {
  account_id = "%[2]s"
  mode       = "exclude"

  tunnels {
    address     = "192.168.0.0/16"
    description = "Private network"
  }

  tunnels {
    host        = "internal.local"
    description = "Internal domain"
  }
}

resource "cloudflare_split_tunnel" "%[1]s_include" {
  account_id = "%[2]s"
  mode       = "include"

  tunnels {
    address     = "203.0.113.0/24"
    description = "Corporate VPN"
  }
}`, rnd, accountID)

	stateChecks := []statecheck.StateCheck{
		// Verify default profile exists
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("id"), knownvalue.NotNull()),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),

		// Verify both exclude and include arrays are populated
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("exclude"), knownvalue.SetSizeExact(2)),
		statecheck.ExpectKnownValue("cloudflare_zero_trust_device_default_profile."+rnd, tfjsonpath.New("include"), knownvalue.SetSizeExact(1)),
	}

	migrationSteps := acctest.MigrationV2TestStepWithStateNormalization(t, v4Config, tmpDir, "4.52.1", "v4", "v5", stateChecks)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
			acctest.TestAccPreCheck_AccountID(t)
		},
		WorkingDir: tmpDir,
		Steps: append([]resource.TestStep{
			{
				// Step 1: Create with v4 provider
				ExternalProviders: map[string]resource.ExternalProvider{
					"cloudflare": {
						Source:            "cloudflare/cloudflare",
						VersionConstraint: "4.52.1",
					},
				},
				Config: v4Config,
			},
		}, migrationSteps...),
	})
}
