package run

// Status to report to runner
type Status string

// Statuses for runner
const (
	Info    Status = "info"
	Error   Status = "error"
	Success Status = "success"
	Failed  Status = "failed"
)

// Send message with certain status
func (status Status) Send(cmd *Cmd, msg *Msg) error {
	msg.Status = status
	return msg.Send(cmd.Namespace)
}
