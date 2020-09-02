package runner

// Status run status
type Status string

// ToString convert message to string with current status
func (status Status) ToString(message *Message) (string, error) {
	message.Status = status
	return message.ToString()
}

// Send message with certain status
func (status Status) Send(namespace string, Message *Message) error {
	Message.Status = status
	return Message.Send(namespace)
}
