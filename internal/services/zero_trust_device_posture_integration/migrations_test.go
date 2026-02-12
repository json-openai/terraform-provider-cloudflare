package zero_trust_device_posture_integration_test

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

// TestMigrateDevicePostureIntegrationDeprecatedName tests migration from deprecated resource name
// Uses CrowdStrike S2S integration type which requires:
// - CLOUDFLARE_CROWDSTRIKE_CLIENT_ID
// - CLOUDFLARE_CROWDSTRIKE_CLIENT_SECRET
// - CLOUDFLARE_CROWDSTRIKE_API_URL
// - CLOUDFLARE_CROWDSTRIKE_CUSTOMER_ID
func TestMigrateDevicePostureIntegrationDeprecatedName(t *testing.T) {
	// Zero Trust resources require API_KEY + EMAIL, not API_TOKEN
	originalToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if originalToken != "" {
		os.Unsetenv("CLOUDFLARE_API_TOKEN")
		defer os.Setenv("CLOUDFLARE_API_TOKEN", originalToken)
	}

	clientID := os.Getenv("CLOUDFLARE_CROWDSTRIKE_CLIENT_ID")
	clientSecret := os.Getenv("CLOUDFLARE_CROWDSTRIKE_CLIENT_SECRET")
	apiURL := os.Getenv("CLOUDFLARE_CROWDSTRIKE_API_URL")
	customerID := os.Getenv("CLOUDFLARE_CROWDSTRIKE_CUSTOMER_ID")

	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	rnd := utils.GenerateRandomResourceName()
	name := fmt.Sprintf("tf-migrate-deprecated-%s", rnd)
	tmpDir := t.TempDir()

	// Use DEPRECATED v4 resource name: cloudflare_device_posture_integration
	// Migration should rename to: cloudflare_zero_trust_device_posture_integration
	v4Config := fmt.Sprintf(`
resource "cloudflare_device_posture_integration" "%[1]s" {
  account_id = "%[2]s"
  name       = "%[3]s"
  type       = "crowdstrike_s2s"
  interval   = "24h"

  config {
    api_url       = "%[4]s"
    client_id     = "%[5]s"
    client_secret = "%[6]s"
    customer_id   = "%[7]s"
  }
}`, rnd, accountID, name, apiURL, clientID, clientSecret, customerID)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
			acctest.TestAccPreCheck_AccountID(t)
			acctest.TestAccPreCheck_CrowdStrike(t)
		},
		WorkingDir: tmpDir,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with v4 provider using DEPRECATED name
				ExternalProviders: map[string]resource.ExternalProvider{
					"cloudflare": {
						Source:            "cloudflare/cloudflare",
						VersionConstraint: "4.52.1",
					},
				},
				Config: v4Config,
			},
			// Step 2: Migrate and verify resource name changed
			acctest.MigrationV2TestStep(t, v4Config, tmpDir, "4.52.1", "v4", "v5", []statecheck.StateCheck{
				// Verify resource was renamed from cloudflare_device_posture_integration to cloudflare_zero_trust_device_posture_integration
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("name"), knownvalue.StringExact(name)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("type"), knownvalue.StringExact("crowdstrike_s2s")),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("interval"), knownvalue.StringExact("24h")),

				// Config should be converted from block syntax to attribute syntax
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("config").AtMapKey("api_url"), knownvalue.StringExact(apiURL)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("config").AtMapKey("client_id"), knownvalue.StringExact(clientID)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("config").AtMapKey("customer_id"), knownvalue.StringExact(customerID)),
			}),
		},
	})
}

// TestMigrateDevicePostureIntegrationCurrentName tests migration from current (non-deprecated) resource name
// Uses CrowdStrike S2S integration type
func TestMigrateDevicePostureIntegrationCurrentName(t *testing.T) {
	// Zero Trust resources require API_KEY + EMAIL, not API_TOKEN
	originalToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if originalToken != "" {
		os.Unsetenv("CLOUDFLARE_API_TOKEN")
		defer os.Setenv("CLOUDFLARE_API_TOKEN", originalToken)
	}

	clientID := os.Getenv("CLOUDFLARE_CROWDSTRIKE_CLIENT_ID")
	clientSecret := os.Getenv("CLOUDFLARE_CROWDSTRIKE_CLIENT_SECRET")
	apiURL := os.Getenv("CLOUDFLARE_CROWDSTRIKE_API_URL")
	customerID := os.Getenv("CLOUDFLARE_CROWDSTRIKE_CUSTOMER_ID")

	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	rnd := utils.GenerateRandomResourceName()
	name := fmt.Sprintf("tf-migrate-current-%s", rnd)
	tmpDir := t.TempDir()

	// Use CURRENT v4 resource name: cloudflare_zero_trust_device_posture_integration
	// Migration should keep the same name
	v4Config := fmt.Sprintf(`
resource "cloudflare_zero_trust_device_posture_integration" "%[1]s" {
  account_id = "%[2]s"
  name       = "%[3]s"
  type       = "crowdstrike_s2s"
  interval   = "24h"

  config {
    api_url       = "%[4]s"
    client_id     = "%[5]s"
    client_secret = "%[6]s"
    customer_id   = "%[7]s"
  }
}`, rnd, accountID, name, apiURL, clientID, clientSecret, customerID)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
			acctest.TestAccPreCheck_AccountID(t)
			acctest.TestAccPreCheck_CrowdStrike(t)
		},
		WorkingDir: tmpDir,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with v4 provider using CURRENT name
				ExternalProviders: map[string]resource.ExternalProvider{
					"cloudflare": {
						Source:            "cloudflare/cloudflare",
						VersionConstraint: "4.52.1",
					},
				},
				Config: v4Config,
			},
			// Step 2: Migrate and verify config block → attribute transformation
			acctest.MigrationV2TestStep(t, v4Config, tmpDir, "4.52.1", "v4", "v5", []statecheck.StateCheck{
				// Resource name should stay the same (already using current name)
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("name"), knownvalue.StringExact(name)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("type"), knownvalue.StringExact("crowdstrike_s2s")),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("interval"), knownvalue.StringExact("24h")),

				// Config should be converted from block syntax to attribute syntax
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("config").AtMapKey("api_url"), knownvalue.StringExact(apiURL)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("config").AtMapKey("client_id"), knownvalue.StringExact(clientID)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("config").AtMapKey("customer_id"), knownvalue.StringExact(customerID)),
			}),
		},
	})
}

// TestMigrateDevicePostureIntegrationWithIdentifier tests migration when identifier field exists
// Migration should remove the deprecated identifier field
// Uses CrowdStrike S2S integration type
func TestMigrateDevicePostureIntegrationWithIdentifier(t *testing.T) {
	// Zero Trust resources require API_KEY + EMAIL, not API_TOKEN
	originalToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if originalToken != "" {
		os.Unsetenv("CLOUDFLARE_API_TOKEN")
		defer os.Setenv("CLOUDFLARE_API_TOKEN", originalToken)
	}

	clientID := os.Getenv("CLOUDFLARE_CROWDSTRIKE_CLIENT_ID")
	clientSecret := os.Getenv("CLOUDFLARE_CROWDSTRIKE_CLIENT_SECRET")
	apiURL := os.Getenv("CLOUDFLARE_CROWDSTRIKE_API_URL")
	customerID := os.Getenv("CLOUDFLARE_CROWDSTRIKE_CUSTOMER_ID")

	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	rnd := utils.GenerateRandomResourceName()
	name := fmt.Sprintf("tf-migrate-identifier-%s", rnd)
	tmpDir := t.TempDir()

	// V4 config WITH deprecated identifier field (removed in v5)
	v4Config := fmt.Sprintf(`
resource "cloudflare_device_posture_integration" "%[1]s" {
  account_id = "%[2]s"
  name       = "%[3]s"
  type       = "crowdstrike_s2s"
  interval   = "24h"
  identifier = "legacy-identifier-123"

  config {
    api_url       = "%[4]s"
    client_id     = "%[5]s"
    client_secret = "%[6]s"
    customer_id   = "%[7]s"
  }
}`, rnd, accountID, name, apiURL, clientID, clientSecret, customerID)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
			acctest.TestAccPreCheck_AccountID(t)
			acctest.TestAccPreCheck_CrowdStrike(t)
		},
		WorkingDir: tmpDir,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with v4 provider WITH identifier
				ExternalProviders: map[string]resource.ExternalProvider{
					"cloudflare": {
						Source:            "cloudflare/cloudflare",
						VersionConstraint: "4.52.1",
					},
				},
				Config: v4Config,
			},
			// Step 2: Migrate and verify identifier field is removed
			acctest.MigrationV2TestStep(t, v4Config, tmpDir, "4.52.1", "v4", "v5", []statecheck.StateCheck{
				// Verify resource renamed and identifier removed
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("name"), knownvalue.StringExact(name)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("type"), knownvalue.StringExact("crowdstrike_s2s")),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("interval"), knownvalue.StringExact("24h")),

				// Config should be converted from block to attribute
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("config").AtMapKey("api_url"), knownvalue.StringExact(apiURL)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("config").AtMapKey("client_id"), knownvalue.StringExact(clientID)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd, tfjsonpath.New("config").AtMapKey("customer_id"), knownvalue.StringExact(customerID)),
				// Note: identifier field should be removed from state (no check for it)
			}),
		},
	})
}

// TestMigrateDevicePostureIntegrationComprehensive is a comprehensive test covering all migration scenarios:
// - Deprecated resource name rename
// - Current resource name (no rename)
// - Different interval values
// - Identifier field removal
// - Config block → attribute conversion
// - Empty config fields → null transformation
// - Multiple resources in single migration run
func TestMigrateDevicePostureIntegrationComprehensive(t *testing.T) {
	// Zero Trust resources require API_KEY + EMAIL, not API_TOKEN
	originalToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	if originalToken != "" {
		os.Unsetenv("CLOUDFLARE_API_TOKEN")
		defer os.Setenv("CLOUDFLARE_API_TOKEN", originalToken)
	}

	clientID := os.Getenv("CLOUDFLARE_CROWDSTRIKE_CLIENT_ID")
	clientSecret := os.Getenv("CLOUDFLARE_CROWDSTRIKE_CLIENT_SECRET")
	apiURL := os.Getenv("CLOUDFLARE_CROWDSTRIKE_API_URL")
	customerID := os.Getenv("CLOUDFLARE_CROWDSTRIKE_CUSTOMER_ID")

	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	rnd1 := utils.GenerateRandomResourceName()
	rnd2 := utils.GenerateRandomResourceName()
	rnd3 := utils.GenerateRandomResourceName()
	name1 := fmt.Sprintf("tf-migrate-comprehensive-1-%s", rnd1)
	name2 := fmt.Sprintf("tf-migrate-comprehensive-2-%s", rnd2)
	name3 := fmt.Sprintf("tf-migrate-comprehensive-3-%s", rnd3)
	tmpDir := t.TempDir()

	// Create 3 resources covering different scenarios:
	// 1. Deprecated name with identifier field (should rename & remove identifier)
	// 2. Current name with unusual interval (should preserve interval)
	// 3. Deprecated name with standard interval (should rename, preserve interval)
	v4Config := fmt.Sprintf(`
# Resource 1: Deprecated name + identifier field
resource "cloudflare_device_posture_integration" "%[1]s" {
  account_id = "%[4]s"
  name       = "%[5]s"
  type       = "crowdstrike_s2s"
  interval   = "1h"
  identifier = "legacy-identifier-to-remove"

  config {
    api_url       = "%[8]s"
    client_id     = "%[9]s"
    client_secret = "%[10]s"
    customer_id   = "%[11]s"
  }
}

# Resource 2: Current name with non-standard interval
resource "cloudflare_zero_trust_device_posture_integration" "%[2]s" {
  account_id = "%[4]s"
  name       = "%[6]s"
  type       = "crowdstrike_s2s"
  interval   = "6h"

  config {
    api_url       = "%[8]s"
    client_id     = "%[9]s"
    client_secret = "%[10]s"
    customer_id   = "%[11]s"
  }
}

# Resource 3: Deprecated name with standard interval
resource "cloudflare_device_posture_integration" "%[3]s" {
  account_id = "%[4]s"
  name       = "%[7]s"
  type       = "crowdstrike_s2s"
  interval   = "24h"

  config {
    api_url       = "%[8]s"
    client_id     = "%[9]s"
    client_secret = "%[10]s"
    customer_id   = "%[11]s"
  }
}`, rnd1, rnd2, rnd3, accountID, name1, name2, name3, apiURL, clientID, clientSecret, customerID)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
			acctest.TestAccPreCheck_AccountID(t)
			acctest.TestAccPreCheck_CrowdStrike(t)
		},
		WorkingDir: tmpDir,
		Steps: []resource.TestStep{
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
			// Step 2: Migrate and verify all transformations
			acctest.MigrationV2TestStep(t, v4Config, tmpDir, "4.52.1", "v4", "v5", []statecheck.StateCheck{
				// Resource 1: Verify deprecated name was renamed and identifier removed
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd1, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd1, tfjsonpath.New("name"), knownvalue.StringExact(name1)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd1, tfjsonpath.New("type"), knownvalue.StringExact("crowdstrike_s2s")),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd1, tfjsonpath.New("interval"), knownvalue.StringExact("1h")),
				// Verify config is attribute (not block) and has correct values
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd1, tfjsonpath.New("config").AtMapKey("api_url"), knownvalue.StringExact(apiURL)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd1, tfjsonpath.New("config").AtMapKey("client_id"), knownvalue.StringExact(clientID)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd1, tfjsonpath.New("config").AtMapKey("customer_id"), knownvalue.StringExact(customerID)),
				// Note: identifier field should be removed from state (no check for it - absence is verified by no drift)

				// Resource 2: Verify current name preserved and unusual interval maintained
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd2, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd2, tfjsonpath.New("name"), knownvalue.StringExact(name2)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd2, tfjsonpath.New("type"), knownvalue.StringExact("crowdstrike_s2s")),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd2, tfjsonpath.New("interval"), knownvalue.StringExact("6h")),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd2, tfjsonpath.New("config").AtMapKey("api_url"), knownvalue.StringExact(apiURL)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd2, tfjsonpath.New("config").AtMapKey("client_id"), knownvalue.StringExact(clientID)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd2, tfjsonpath.New("config").AtMapKey("customer_id"), knownvalue.StringExact(customerID)),

				// Resource 3: Verify deprecated name renamed and standard interval preserved
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd3, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd3, tfjsonpath.New("name"), knownvalue.StringExact(name3)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd3, tfjsonpath.New("type"), knownvalue.StringExact("crowdstrike_s2s")),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd3, tfjsonpath.New("interval"), knownvalue.StringExact("24h")),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd3, tfjsonpath.New("config").AtMapKey("api_url"), knownvalue.StringExact(apiURL)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd3, tfjsonpath.New("config").AtMapKey("client_id"), knownvalue.StringExact(clientID)),
				statecheck.ExpectKnownValue("cloudflare_zero_trust_device_posture_integration."+rnd3, tfjsonpath.New("config").AtMapKey("customer_id"), knownvalue.StringExact(customerID)),
			}),
		},
	})
}
