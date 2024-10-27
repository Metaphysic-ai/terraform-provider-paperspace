package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomTemplatesDataSource(t *testing.T) {
	// TODO: Create real test
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "paperspace_custom_templates" "test" {}`,
				Check:  resource.ComposeAggregateTestCheckFunc(
				// // Verify number of custom_templates returned
				// resource.TestCheckResourceAttr("data.paperspace_custom_templates.test", "custom_templates.#", "10"),

				// resource.TestCheckResourceAttr("data.paperspace_custom_templates.test", "custom_templates.0.id", "exampleid"),
				// resource.TestCheckResourceAttr("data.paperspace_custom_templates.test", "custom_templates.0.name", "Example Name"),

				// // Verify number of nested attr returned
				// resource.TestCheckResourceAttr("data.paperspace_custom_templates.test", "custom_templates.0.nested_attr.#", "1"),

				// resource.TestCheckResourceAttr("data.paperspace_custom_templates.test", "custom_templates.0.nested_attr.0.id", "exampleid"),
				),
			},
		},
	})
}
