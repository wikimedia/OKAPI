package task

// State inferface to hold the task state
type State interface {
	Get(name string) (string, error)
	GetInt(name string, initial int) int
	GetString(name string, initial string) string
	Set(name string, value interface{}) error
	Exists(name string) bool
	Clear()
}
