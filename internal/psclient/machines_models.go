package psclient

import "time"

// Accelerator represents the structure for each accelerator in the machine.
type Accelerator struct {
	Name   string `json:"name"`
	Memory string `json:"memory"`
	Count  int64  `json:"count"`
}

// Reservation represents the structure for machine reservation details.
type Reservation struct {
	Name       string    `json:"name"`
	ID         string    `json:"id"`
	DtStarted  time.Time `json:"dtStarted"`
	DtFinished time.Time `json:"dtFinished"`
	IsActive   bool      `json:"isActive"`
}

// Machine represents the full structure of the machine data.
type Machine struct {
	ID                     string        `json:"id"`
	Name                   string        `json:"name"`
	State                  string        `json:"state"`
	OS                     string        `json:"os"`
	MachineType            string        `json:"machineType"`
	AgentType              string        `json:"agentType"`
	CPUs                   int64         `json:"cpus"`
	RAM                    string        `json:"ram"`
	StorageTotal           string        `json:"storageTotal"`
	StorageUsed            string        `json:"storageUsed"`
	Accelerators           []Accelerator `json:"accelerators"`
	RegionFull             string        `json:"region"`
	PrivateIP              string        `json:"privateIp"`
	NetworkID              string        `json:"networkId"`
	PublicIP               *string       `json:"publicIp"` // Nullable field
	PublicIPType           string        `json:"publicIpType"`
	AutoShutdownEnabled    bool          `json:"autoShutdownEnabled"`
	AutoShutdownTimeout    *int64        `json:"autoShutdownTimeout"` // Nullable field
	AutoShutdownForce      bool          `json:"autoShutdownForce"`
	AutoSnapshotEnabled    bool          `json:"autoSnapshotEnabled"`
	AutoSnapshotFrequency  *string       `json:"autoSnapshotFrequency"` // Nullable field
	AutoSnapshotSaveCount  *int64        `json:"autoSnapshotSaveCount"` // Nullable field
	UpdatesPending         bool          `json:"updatesPending"`
	RestorePointEnabled    bool          `json:"restorePointEnabled"`
	RestorePointFrequency  *string       `json:"restorePointFrequency"`  // Nullable field
	RestorePointSnapshotID *string       `json:"restorePointSnapshotId"` // Nullable field
	UsageRate              float64       `json:"usageRate"`
	StorageRate            float64       `json:"storageRate"`
	DtCreated              string        `json:"dtCreated"`
	DtModified             string        `json:"dtModified"`
	DtDeleted              *string       `json:"dtDeleted"`   // Nullable field
	Reservation            *Reservation  `json:"reservation"` // Nullable field
}

type MachineCreateConfig struct {
	Name                  string `json:"name"`        // required
	MachineType           string `json:"machineType"` // required
	TemplateID            string `json:"templateId"`  // required
	DiskSize              int64  `json:"diskSize"`    // required
	Region                string `json:"region"`      // required
	NetworkID             string `json:"networkId,omitempty"`
	PublicIPType          string `json:"publicIpType,omitempty"`
	StartOnCreate         bool   `json:"startOnCreate"`
	AutoSnapshotEnabled   *bool  `json:"autoSnapshotEnabled,omitempty"`
	AutoSnapshotFrequency string `json:"autoSnapshotFrequency,omitempty"`
	AutoSnapshotSaveCount *int64 `json:"autoSnapshotSaveCount,omitempty"`
	AutoShutdownEnabled   *bool  `json:"autoShutdownEnabled,omitempty"`
	AutoShutdownTimeout   *int64 `json:"autoShutdownTimeout,omitempty"`
	AutoShutdownForce     *bool  `json:"autoShutdownForce,omitempty"`
	EnableNvlink          *bool  `json:"enableNvlink,omitempty"`
	TakeInitialSnapshot   *bool  `json:"takeInitialSnapshot,omitempty"`
	StartupScriptID       string `json:"startupScriptId,omitempty"`
	EmailPassword         *bool  `json:"emailPassword,omitempty"`

	// TODO: Consider pointer+omit here
	AccessorIDs []string `json:"accessorIds,omitempty"`
}

type MachineUpdateConfig struct {
	Name                  string `json:"name,omitempty"`
	MachineType           string `json:"machineType,omitempty"`
	DiskSize              int64  `json:"diskSize,omitempty"`
	NetworkID             string `json:"networkId,omitempty"`
	PublicIPType          string `json:"publicIpType,omitempty"`
	AutoSnapshotEnabled   *bool  `json:"autoSnapshotEnabled,omitempty"`
	AutoSnapshotFrequency string `json:"autoSnapshotFrequency,omitempty"`
	AutoSnapshotSaveCount *int64 `json:"autoSnapshotSaveCount,omitempty"`
	AutoShutdownEnabled   *bool  `json:"autoShutdownEnabled,omitempty"`
	AutoShutdownTimeout   *int64 `json:"autoShutdownTimeout,omitempty"`
	AutoShutdownForce     *bool  `json:"autoShutdownForce,omitempty"`
}

type Event struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	State      string  `json:"state"`
	MachineID  string  `json:"machineId"`
	DtCreated  string  `json:"dtCreated"`
	DtStarted  *string `json:"dtStarted"`
	DtFinished *string `json:"dtFinished"`
	Error      *string `json:"error"`
}

type MashineResponse struct {
	Event Event   `json:"event"`
	Data  Machine `json:"data"`
}
