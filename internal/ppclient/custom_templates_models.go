package ppclient

type AvailableMachineType struct {
	MachineTypeLabel string `json:"machineTypeLabel"`
	IsAvailable      bool   `json:"isAvailable"`
}

type CustomTemplate struct {
	ID                    string                 `json:"id"`
	Name                  string                 `json:"name"`
	AgentType             string                 `json:"agentType"`
	OperatingSystemLabel  string                 `json:"operatingSystemLabel"`
	Region                string                 `json:"region"`
	DefaultSizeGb         int                    `json:"defaultSizeGb"`
	AvailableMachineTypes []AvailableMachineType `json:"availableMachineTypes"`
	ParentMachineID       string                 `json:"parentMachineId"`
	DtCreated             string                 `json:"dtCreated"`
	DtDeleted             *string                `json:"dtDeleted"` // Nullable
}
