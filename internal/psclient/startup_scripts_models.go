package psclient

// Accelerator represents the structure for each accelerator in the machine.
type StartupScriptCreateConfig struct {
	Name      string `json:"name"`
	Script    string `json:"script"`
	IsRunOnce bool   `json:"isRunOnce"`
}

type StartupScript struct {
	ID                 string   `json:"id"`                 // The ID of the startup script
	Name               string   `json:"name"`               // The name of the startup script
	Description        *string  `json:"description"`        // The description of the startup script (nullable)
	IsEnabled          bool     `json:"isEnabled"`          // Whether the startup script is enabled
	IsRunOnce          bool     `json:"isRunOnce"`          // Whether the startup script runs once or on every boot
	AssignedMachineIDs []string `json:"assignedMachineIds"` // The IDs of machines assigned to this script
	DtCreated          string   `json:"dtCreated"`          // The creation date of the startup script
	DtDeleted          *string  `json:"dtDeleted"`          // The deletion date of the startup script (nullable)
}
