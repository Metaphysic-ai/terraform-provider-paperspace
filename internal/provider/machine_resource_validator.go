package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *machineResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data machineResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Auto Snapshot

	// If attribute is null, unknown or false
	if !data.AutoSnapshotEnabled.ValueBool() {
		if !data.AutoSnapshotFrequency.IsNull() {
			addAttributeDepsError(&resp.Diagnostics, "auto_snapshot_frequency", []string{"auto_snapshot_enabled"})
		}

		if !data.AutoSnapshotSaveCount.IsNull() {
			addAttributeDepsError(&resp.Diagnostics, "auto_snapshot_save_count", []string{"auto_snapshot_enabled"})
		}
	} else {
		if data.AutoSnapshotFrequency.IsNull() {
			addAttributeDepsError(&resp.Diagnostics, "auto_snapshot_enabled", []string{"auto_snapshot_frequency"})
		}

		if data.AutoSnapshotSaveCount.IsNull() {
			addAttributeDepsError(&resp.Diagnostics, "auto_snapshot_enabled", []string{"auto_snapshot_save_count"})
		}
	}

	// Auto Shutdown

	// If attribute is null, unknown or false
	if !data.AutoShutdownEnabled.ValueBool() {
		if !data.AutoShutdownTimeout.IsNull() {
			addAttributeDepsError(&resp.Diagnostics, "auto_shutdown_timeout", []string{"auto_shutdown_enabled"})
		}

		if !data.AutoShutdownForce.IsNull() {
			addAttributeDepsError(&resp.Diagnostics, "auto_shutdown_force", []string{"auto_shutdown_enabled"})
		}
	} else { // Attribute is true
		if data.AutoShutdownTimeout.IsNull() {
			addAttributeDepsError(&resp.Diagnostics, "auto_shutdown_enabled", []string{"auto_shutdown_timeout"})
		}
	}

	// Restore Point

	// NOTE: Restore point feature is not implemented yet
	// If attribute is null, unknown or false
	if !data.RestorePointEnabled.ValueBool() {
		if !data.RestorePointFrequency.IsNull() {
			addAttributeDepsError(&resp.Diagnostics, "restore_point_frequency", []string{"restore_point_enabled"})
		}
	} else { // Attribute is true
		if data.RestorePointFrequency.IsNull() {
			addAttributeDepsError(&resp.Diagnostics, "restore_point_enabled", []string{"restore_point_frequency"})
		}
	}
}

func addAttributeDepsError(diags *diag.Diagnostics, attrName string, deps []string) {
	errorType := "Missing Attribute Configuration"
	depsString := strings.Join(deps, ", ")

	diags.AddAttributeError(
		path.Root(attrName),
		errorType,
		fmt.Sprintf("Attribute %s requires also: %s", attrName, depsString),
	)
}
