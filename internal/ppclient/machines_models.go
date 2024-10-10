package ppclient

import "time"

// Accelerator represents the structure for each accelerator in the machine.
type Accelerator struct {
	Name   string `json:"name"`
	Memory string `json:"memory"`
	Count  int    `json:"count"`
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
	CPUs                   int           `json:"cpus"`
	RAM                    string        `json:"ram"`
	StorageTotal           string        `json:"storageTotal"`
	StorageUsed            string        `json:"storageUsed"`
	Accelerators           []Accelerator `json:"accelerators"`
	Region                 string        `json:"region"`
	PrivateIP              string        `json:"privateIp"`
	NetworkID              string        `json:"networkId"`
	PublicIP               *string       `json:"publicIp"` // Nullable field
	PublicIPType           string        `json:"publicIpType"`
	AutoShutdownEnabled    bool          `json:"autoShutdownEnabled"`
	AutoShutdownTimeout    *int          `json:"autoShutdownTimeout"` // Nullable field
	AutoShutdownForce      bool          `json:"autoShutdownForce"`
	AutoSnapshotEnabled    bool          `json:"autoSnapshotEnabled"`
	AutoSnapshotFrequency  *string       `json:"autoSnapshotFrequency"` // Nullable field
	AutoSnapshotSaveCount  *int          `json:"autoSnapshotSaveCount"` // Nullable field
	UpdatesPending         bool          `json:"updatesPending"`
	RestorePointEnabled    bool          `json:"restorePointEnabled"`
	RestorePointFrequency  *string       `json:"restorePointFrequency"`  // Nullable field
	RestorePointSnapshotID *string       `json:"restorePointSnapshotId"` // Nullable field
	UsageRate              float64       `json:"usageRate"`
	StorageRate            float64       `json:"storageRate"`
	DtCreated              time.Time     `json:"dtCreated"`
	DtModified             time.Time     `json:"dtModified"`
	DtDeleted              *time.Time    `json:"dtDeleted"`   // Nullable field
	Reservation            *Reservation  `json:"reservation"` // Nullable field
}

type MashinesResponse struct {
	HasMore  bool      `json:"hasMore"`
	NextPage string    `json:"nextPage"`
	Items    []Machine `json:"items"`
}

// Configuration of machine to be created
type MachineConfig struct {
	Name          string `json:"name"`        // required
	MachineType   string `json:"machineType"` // required
	TemplateID    string `json:"templateId"`  // required
	DiskSize      int    `json:"diskSize"`    // required
	Region        string `json:"region"`      // required
	NetworkID     string `json:"networkId"`
	PublicIpType  string `json:"publicIpType"`
	StartOnCreate bool   `json:"startOnCreate"`

	// TODO: Add remaining attributes

	// AutoSnapshotEnabled bool   `json:"autoSnapshotEnabled"`
	// AutoSnapshotFrequency string `json:"autoSnapshotFrequency"`
	// AutoSnapshotSaveCount int  `json:"autoSnapshotSaveCount"`

	// AutoShutdownEnabled   bool `json:"autoShutdownEnabled"`
	// AutoShutdownTimeout   int  `json:"autoShutdownTimeout"`
	// AutoShutdownForce     bool `json:"autoShutdownForce"`

	// RestorePointEnabled   bool `json:"restorePointEnabled"`
	// RestorePointFrequency string `json:"restorePointFrequency"`

	// EnableNvlink    bool   `json:"enableNvlink"`
	// TakeInitialSnapshot   bool `json:"takeInitialSnapshot"`

	// // Useful
	// StartupScriptID string   `json:"startupScriptId"`
	// EmailPassword   bool     `json:"emailPassword"` # default: true
	AccessorIds []string `json:"accessorIds"`
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
