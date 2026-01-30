package leaked_credential_check_rule_test

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

// TestMigrateLeakedCredentialCheckRule_V4ToV5_BasicRule tests migration of a basic rule with all fields
func TestMigrateLeakedCredentialCheckRule_V4ToV5_BasicRule(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	resourceName := "cloudflare_leaked_credential_check_rule." + rnd
	tmpDir := t.TempDir()
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")

	// Create v4 configuration with all fields
	v4Config := fmt.Sprintf(`
resource "cloudflare_leaked_credential_check_rule" "%[1]s" {
  zone_id  = "%[2]s"
  username = "lookup_json_string(http.request.body.raw, \"user\")"
  password = "lookup_json_string(http.request.body.raw, \"pass\")"
}`, rnd, zoneID)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
			acctest.TestAccPreCheck_ZoneID(t)
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "zone_id", zoneID),
					resource.TestCheckResourceAttr(resourceName, "username", "lookup_json_string(http.request.body.raw, \"user\")"),
					resource.TestCheckResourceAttr(resourceName, "password", "lookup_json_string(http.request.body.raw, \"pass\")"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Step 2: Run migration and verify state
			// Note: Resource name stays the same (cloudflare_leaked_credential_check_rule)
			acctest.MigrationV2TestStep(t, v4Config, tmpDir, "4.52.1", "v4", "v5", []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					resourceName,
					tfjsonpath.New("zone_id"),
					knownvalue.StringExact(zoneID),
				),
				statecheck.ExpectKnownValue(
					resourceName,
					tfjsonpath.New("username"),
					knownvalue.StringExact("lookup_json_string(http.request.body.raw, \"user\")"),
				),
				statecheck.ExpectKnownValue(
					resourceName,
					tfjsonpath.New("password"),
					knownvalue.StringExact("lookup_json_string(http.request.body.raw, \"pass\")"),
				),
			}),
		},
	})
}

// TestMigrateLeakedCredentialCheckRule_V4ToV5_ComplexExpressions tests migration with complex field expressions
func TestMigrateLeakedCredentialCheckRule_V4ToV5_ComplexExpressions(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	resourceName := "cloudflare_leaked_credential_check_rule." + rnd
	tmpDir := t.TempDir()
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")

	// Create v4 configuration with complex expressions
	v4Config := fmt.Sprintf(`
resource "cloudflare_leaked_credential_check_rule" "%[1]s" {
  zone_id  = "%[2]s"
  username = "lookup_json_string(lookup_json_string(http.request.body.raw, \"payload\"), \"username\")"
  password = "http.request.headers[\"x-password\"][0]"
}`, rnd, zoneID)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.TestAccPreCheck(t)
			acctest.TestAccPreCheck_ZoneID(t)
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "zone_id", zoneID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Step 2: Run migration and verify state
			acctest.MigrationV2TestStep(t, v4Config, tmpDir, "4.52.1", "v4", "v5", []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					resourceName,
					tfjsonpath.New("zone_id"),
					knownvalue.StringExact(zoneID),
				),
				statecheck.ExpectKnownValue(
					resourceName,
					tfjsonpath.New("username"),
					knownvalue.StringExact("lookup_json_string(lookup_json_string(http.request.body.raw, \"payload\"), \"username\")"),
				),
				statecheck.ExpectKnownValue(
					resourceName,
					tfjsonpath.New("password"),
					knownvalue.StringExact("http.request.headers[\"x-password\"][0]"),
				),
			}),
		},
	})
}
