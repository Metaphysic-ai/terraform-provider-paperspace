package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccMachineResourceName = "paperspace_machine.test"

// TODO: Consider using the files for configurations, it's natively supported by []resource.TestStep
// TODO: Add 'startup_script_id' to test (once provider is able to create startup scripts)

var testAccMachineResourceConfigs = map[string]string{
	"CreateRead": providerConfig + `
resource "paperspace_machine" "test" {
  name            = "paperspace-provider-test-CreateRead"
  machine_type    = "C2"
  template_id  	  = "t0nspur5"
  disk_size       = 50
  region          = "ny2"

  # Only for new machines
  enable_nvlink 		= false
  take_initial_snapshot = true
  email_password 		= false
}
`,

	"UpdateRead": providerConfig + `
resource "paperspace_machine" "test" {
  name            = "paperspace-provider-test-UpdateRead"
  machine_type    = "C3"
  template_id  	  = "t0nspur5"
  disk_size       = 100
  region          = "ny2"
  public_ip_type  = "static"
  accessor_ids    = []
  state 		  = "off"

  auto_snapshot_enabled = true
  auto_snapshot_save_count = 1
  auto_snapshot_frequency  = "daily"

  auto_shutdown_enabled = true
  auto_shutdown_force   = true
  auto_shutdown_timeout = 1
}
`,

	// Only required fields are set here, to  test defaults
	"CreateReadDefaults": providerConfig + `
resource "paperspace_machine" "test_defaults" {
  name            = "paperspace-provider-test-CreateReadDefaults"
  machine_type    = "C2"
  template_id  	  = "t0nspur5"
  disk_size       = 50
  region          = "ny2"
}
`,

	"CreateStart": providerConfig + `
resource "paperspace_machine" "test_start" {
  name            = "paperspace-provider-test-CreateStart"
  machine_type    = "C2"
  template_id  	  = "t0nspur5"
  disk_size       = 50
  region          = "ny2"

  state 		  = "ready"
  email_password  = false
}
`,

	"UpdateStarted": providerConfig + `
resource "paperspace_machine" "test_start" {
  name            = "paperspace-provider-test-UpdateStarted"
  machine_type    = "C2"
  template_id  	  = "t0nspur5"
  disk_size       = 100
  region          = "ny2"

  state 		  = "ready"
  email_password  = false
}
`,

	"StopStarted": providerConfig + `
resource "paperspace_machine" "test_start" {
  name            = "paperspace-provider-test-UpdateStarted"
  machine_type    = "C2"
  template_id  	  = "t0nspur5"
  disk_size       = 100
  region          = "ny2"

  state 		  = "off"
  email_password  = false
}
`,
}

func TestAccMachineResource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMachineResourceConfigs["CreateRead"],
				Check: resource.ComposeAggregateTestCheckFunc(genTestCheckFuncs(
					testAccMachineResourceName,
					map[string]string{
						"accessor_ids.#": "0",

						// Verify created machine attributes
						"name":                  "paperspace-provider-test-CreateRead",
						"machine_type":          "C2",
						"template_id":           "t0nspur5",
						"disk_size":             "50",
						"region":                "ny2",
						"public_ip_type":        "dynamic",
						"enable_nvlink":         "false",
						"take_initial_snapshot": "true",
						"email_password":        "false",
						"state":                 "off",

						// Verify machine has Computed attributes filled
						"cpus":        "1",
						"region_full": "East Coast (NY2)",

						// Verify dynamic values have any value set in the state
						"id":            "_any_",
						"private_ip":    "_any_",
						"ram":           "_any_",
						"storage_total": "_any_",
						"storage_used":  "_any_",
						"usage_rate":    "_any_",
						"storage_rate":  "_any_",
					},
				)...),
			},
			// TODO: Add this once import is implemented
			// ImportState testing
			// {
			// 	ResourceName:      testAccMachineResourceName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// 	// The last_updated attribute does not exist in the Paperspace
			// 	// API, therefore there is no value for it during import.
			// 	ImportStateVerifyIgnore: []string{"last_updated"},
			// },

			// Update and Read testing
			{
				Config: testAccMachineResourceConfigs["UpdateRead"],
				Check: resource.ComposeAggregateTestCheckFunc(genTestCheckFuncs(
					testAccMachineResourceName,
					map[string]string{
						"accessor_ids.#": "0",

						// Verify machine updated
						"name":                     "paperspace-provider-test-UpdateRead",
						"machine_type":             "C3",
						"template_id":              "t0nspur5",
						"disk_size":                "100",
						"region":                   "ny2",
						"public_ip_type":           "static",
						"state":                    "off",
						"auto_snapshot_enabled":    "true",
						"auto_snapshot_save_count": "1",
						"auto_snapshot_frequency":  "daily",
						"auto_shutdown_enabled":    "true",
						"auto_shutdown_force":      "true",
						"auto_shutdown_timeout":    "1",
						"cpus":                     "2",
						"region_full":              "East Coast (NY2)",

						// Verify dynamic values have any value set in the state.
						"id":            "_any_",
						"public_ip":     "_any_",
						"ram":           "_any_",
						"storage_total": "_any_",
						"storage_used":  "_any_",
						"usage_rate":    "_any_",
						"storage_rate":  "_any_",
					},
				)...),
			},
		},
	})
}

// Test Defaults

func TestAccMachineResourceDefaults(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMachineResourceConfigs["CreateReadDefaults"],
				Check: resource.ComposeAggregateTestCheckFunc(genTestCheckFuncs(
					"paperspace_machine.test_defaults",
					map[string]string{
						"auto_shutdown_enabled": "false",
						"auto_shutdown_force":   "false",
						"auto_snapshot_enabled": "false",
						"email_password":        "true",
						"enable_nvlink":         "false",
						"public_ip_type":        "dynamic",
						"restore_point_enabled": "false",
						"state":                 "off",
						"take_initial_snapshot": "false",

						"auto_shutdown_timeout":    "null",
						"auto_snapshot_frequency":  "null",
						"auto_snapshot_save_count": "null",
						"restore_point_frequency":  "null",
						"startup_script_id":        "null",
					},
				)...),
			},
		},
	})
}

// Test Create and Start.
func TestAccMachineResourceCreateStartUpdateStop(t *testing.T) {
	// Especially useful for CI, to skip test using 'go test -short'
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMachineResourceConfigs["CreateStart"],
				Check: resource.ComposeAggregateTestCheckFunc(genTestCheckFuncs(
					"paperspace_machine.test_start",
					map[string]string{
						"state": "ready",
					},
				)...),
			},
			{
				Config: testAccMachineResourceConfigs["UpdateStarted"],
				Check: resource.ComposeAggregateTestCheckFunc(genTestCheckFuncs(
					"paperspace_machine.test_start",
					map[string]string{
						"name":      "paperspace-provider-test-UpdateStarted",
						"state":     "ready",
						"disk_size": "100",
					},
				)...),
			},
			{
				Config: testAccMachineResourceConfigs["StopStarted"],
				Check: resource.ComposeAggregateTestCheckFunc(genTestCheckFuncs(
					"paperspace_machine.test_start",
					map[string]string{
						"state": "off",
					},
				)...),
			},
		},
	})
}

// Private

func genTestCheckFuncs(resourceName string, attributes map[string]string) []resource.TestCheckFunc {
	var checks []resource.TestCheckFunc

	for key, value := range attributes {
		if value == "null" {
			checks = append(checks, resource.TestCheckNoResourceAttr(resourceName, key))
		} else if value == "_any_" {
			checks = append(checks, resource.TestCheckResourceAttrSet(resourceName, key))
		} else {
			checks = append(checks, resource.TestCheckResourceAttr(resourceName, key, value))
		}
	}

	return checks
}
