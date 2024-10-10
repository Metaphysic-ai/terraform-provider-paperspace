package ppclient

type Team struct {
	Name                  string  `json:"name"`
	ID                    string  `json:"id"`
	Namespace             string  `json:"namespace"`
	PublicProfileImageUrl *string `json:"publicProfileImageUrl"`
	IsUserTeam            bool    `json:"isUserTeam"`
	DtCreated             string  `json:"dtCreated"`
}

type TeamMembership struct {
	Team    Team `json:"team"`
	IsOwner bool `json:"isOwner"`
	IsAdmin bool `json:"isAdmin"`
}

type Metadata struct {
	Tags string `json:"tags"`
}

type User struct {
	FirstName                 string           `json:"firstName"`
	LastName                  string           `json:"lastName"`
	Email                     string           `json:"email"`
	DtCreated                 string           `json:"dtCreated"`
	DtConfirmed               string           `json:"dtConfirmed"`
	TeamMemberships           []TeamMembership `json:"teamMemberships"`
	IsPhoneVerified           bool             `json:"isPhoneVerified"`
	IsPasswordAuthEnabled     bool             `json:"isPasswordAuthEnabled"`
	IsQrCodeBasedMfaEnabled   bool             `json:"isQrCodeBasedMfaEnabled"`
	IsQrCodeBasedMfaConfirmed bool             `json:"isQrCodeBasedMfaConfirmed"`
	Preferences               *string          `json:"preferences"`
	ID                        string           `json:"id"`
	Metadata                  Metadata         `json:"metadata"`
}

type TeamInfo struct {
	Namespace   string `json:"namespace"`
	IsPrivate   bool   `json:"isPrivate"`
	MaxMachines int    `json:"maxMachines"`
	ID          string `json:"id"`
}

type AuthSession struct {
	User User     `json:"user"`
	Team TeamInfo `json:"team"`
}
