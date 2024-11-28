package psclient

type PrivateNetwork struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Region    string  `json:"region"`
	Network   string  `json:"network"`
	Netmask   string  `json:"netmask"`
	DtCreated string  `json:"dtCreated"`
	DtDeleted *string `json:"dtDeleted"` // Nullable
}
