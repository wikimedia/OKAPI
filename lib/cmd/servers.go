package cmd

// Server server call type
type Server string

// Available servers
const (
	Runner Server = "runner"
	Stream Server = "stream"
	Queue  Server = "queue"
	API    Server = "api"
)
