---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "paperspace_machine Resource - paperspace"
subcategory: ""
description: |-
  Machine resource
---

# paperspace_machine (Resource)

Machine resource

## Example Usage

```terraform
# Manage example machine
resource "paperspace_machine" "example" {
  name         = "Example Name"
  machine_type = "C1"
  template_id  = "tkni3aa4"
  disk_size    = 50
  region       = "ny2"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `disk_size` (Number) The disk size in gigabytes. Updates to this field will trigger a stop/start of the machine.
- `machine_type` (String) The machine type. Updates to this field will trigger a stop/start of the machine.
- `name` (String) The name of the new machine.
- `region` (String) The region to create the machine in.
- `template_id` (String) The template ID.

### Optional

- `accessor_ids` (List of String) The IDs of users to grant access to the machine. Applies only on resource creation.
- `auto_shutdown_enabled` (Boolean) Whether to enable auto shutdown.
- `auto_shutdown_force` (Boolean) Whether to force shutdown the machine. May be troubles with updating the value, seems like Paperspace API issue.Disable auto shutdown and then enable with different option to update.
- `auto_shutdown_timeout` (Number) The auto shutdown timeout in hours. Must be set if `auto_shutdown_enabled` is true. May be troubles with updating the value, seems like Paperspace API issue.Disable auto shutdown and then enable with different option to update.
- `auto_snapshot_enabled` (Boolean) Whether to enable auto snapshots.
- `auto_snapshot_frequency` (String) The auto snapshot frequency. Possible values: `hourly`, `daily`, `weekly`, `monthly`.
- `auto_snapshot_save_count` (Number) The number of auto snapshots to save. Must be between 1 and 9 if `auto_snapshot_enabled` is true.
- `email_password` (Boolean) Whether to email the password. Applies only on resource creation.
- `enable_nvlink` (Boolean) Whether to enable NVLink.
- `private_network_id` (String) Private network ID. You can migrate machines between private networks and from the default network to a private network. It is not possible to migrate a machine back to the default network. If this is required, please file a support ticket.
- `public_ip_type` (String) The public IP type. Possible values: `static`, `dynamic`, `none`.
- `startup_script_id` (String) The startup script ID. Forces resource replacement if changed.
- `state` (String) Desired state of the machine. Possible values: `off`, `ready`.
- `take_initial_snapshot` (Boolean) Whether to take an initial snapshot. Applies only on resource creation.

### Read-Only

- `agent_type` (String) Agent type of the machine.
- `cpus` (Number) Number of CPUs.
- `dt_created` (String) Created date timestamp of the machine.
- `dt_modified` (String) Modified date timestamp of the machine.
- `id` (String) The ID of the machine.
- `os` (String) Operating system of the machine.
- `private_ip` (String) Private IP address of the machine.
- `public_ip` (String) Public IP address of the machine.
- `ram` (String) RAM amount of the machine.
- `region_full` (String) Full machine region name.
- `restore_point_enabled` (Boolean) Whether to use initial snapshot as a restore point.
- `restore_point_frequency` (String) The restore point frequency. Possible values: `shutdown`.
- `restore_point_snapshot_id` (String) The restore point snapshot ID.
- `storage_rate` (Number) Storage rate of the machine.
- `storage_total` (String) Storage total of the machine.
- `storage_used` (String) Storage used of the machine.
- `usage_rate` (Number) Usage rate of the machine.

## Import

Import is supported using the following syntax:

```shell
# Machine can be imported by specifying the identifier.
terraform import paperspace_machine.example ps0123456789
```
