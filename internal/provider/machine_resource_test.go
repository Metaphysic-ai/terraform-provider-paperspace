package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccMachineResourceName = "paperspace_machine.test"

var testAccMachineResourceConfigs = map[string]string{
	"CreateRead": providerConfig + `
resource "paperspace_machine" "test" {
  name            = "paperspace-provider-test-CreateRead"
  machine_type    = "C1"
  template_id  	  = "tkni3aa4"
  disk_size       = 50
  region          = "ny2"
  public_ip_type  = "dynamic"
  start_on_create = false
}
`,

	"UpdateRead": providerConfig + `
resource "paperspace_machine" "test" {
  name            = "paperspace-provider-test-UpdateRead"
  machine_type    = "C3"
  template_id  	  = "tkni3aa4"
  disk_size       = 100
  region          = "ny2"
  public_ip_type  = "static"
  start_on_create = false
}
`,
}

func TestAccMachineResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMachineResourceConfigs["CreateRead"],
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items.
					resource.TestCheckResourceAttr(testAccMachineResourceName, "accessorIds.#", "2"),

					// Verify created machine attributes
					resource.TestCheckResourceAttr(testAccMachineResourceName, "name", "paperspace-provider-test-CreateRead"),
					resource.TestCheckResourceAttr(testAccMachineResourceName, "machine_type", "C1"),
					resource.TestCheckResourceAttr(testAccMachineResourceName, "template_id", "tkni3aa4"),
					resource.TestCheckResourceAttr(testAccMachineResourceName, "disk_size", "50"),
					resource.TestCheckResourceAttr(testAccMachineResourceName, "region", "ny2"),
					resource.TestCheckResourceAttr(testAccMachineResourceName, "public_ip_type", "dynamic"),
					resource.TestCheckResourceAttr(testAccMachineResourceName, "start_on_create", "false"),

					// Verify machine has Computed attributes filled.
					resource.TestCheckResourceAttr(testAccMachineResourceName, "state", "off"),
					resource.TestCheckResourceAttr(testAccMachineResourceName, "cpus", "1"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testAccMachineResourceName, "id"),
					resource.TestCheckResourceAttrSet(testAccMachineResourceName, "private_ip"),
				),
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
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify machine updated.
					resource.TestCheckResourceAttr(testAccMachineResourceName, "name", "paperspace-provider-test-UpdateRead"),
					resource.TestCheckResourceAttr(testAccMachineResourceName, "machine_type", "C3"),
					resource.TestCheckResourceAttr(testAccMachineResourceName, "template_id", "tkni3aa4"),
					resource.TestCheckResourceAttr(testAccMachineResourceName, "disk_size", "100"),
					resource.TestCheckResourceAttr(testAccMachineResourceName, "region", "ny2"),
					resource.TestCheckResourceAttr(testAccMachineResourceName, "public_ip_type", "static"),
					resource.TestCheckResourceAttr(testAccMachineResourceName, "start_on_create", "false"),

					// Verify machine has Computed attributes updated.
					resource.TestCheckResourceAttr(testAccMachineResourceName, "cpus", "2"),

					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testAccMachineResourceName, "id"),
					resource.TestCheckResourceAttrSet(testAccMachineResourceName, "public_ip"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
