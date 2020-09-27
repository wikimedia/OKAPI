package task

// Pool function to get payload into the queue
type Pool func() ([]Payload, error)
