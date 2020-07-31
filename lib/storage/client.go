package storage

var connections = map[ConnectionName]Connection{}

var initializers = map[ConnectionName]func() Connection{
	Local:  localClient,
	Remote: remoteClient,
}

// Available connections
const (
	Local  ConnectionName = "local"
	Remote ConnectionName = "remote"
)

// Client get connection client
func (cName ConnectionName) Client() Connection {
	if con, ok := connections[cName]; ok {
		return con
	}

	connections[cName] = initializers[cName]()

	return connections[cName]
}
