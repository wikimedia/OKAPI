package schedule

// Frequency task execution period
type Frequency string

// All execution timelines
const (
	Once      Frequency = "@once"
	Hourly    Frequency = "@hourly"
	Daily     Frequency = "@daily"
	Weekly    Frequency = "@weekly"
	Monthly   Frequency = "@monthly"
	Quarterly Frequency = "@quarterly"
	Yearly    Frequency = "@yearly"
)

// Info schedule execution context
type Info struct {
	Frequency Frequency `json:"frequency"`
	Workers   int       `json:"workers"`
}
