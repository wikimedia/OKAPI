package roles

// Type roles type identifier
type Type string

// All available roles
const (
	Admin      Type = "admin"
	Client     Type = "client"
	Subscriber Type = "subscriber"
)
