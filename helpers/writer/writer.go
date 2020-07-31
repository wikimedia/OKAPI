package writer

import (
	"archive/tar"
	"io"
	"os"
	"sync"
)

// Payload data to write
type Payload struct {
	ReadCloser io.ReadCloser
	Name       string
	Size       int64
}

// Client writer client struct
type Client struct {
	Queue     chan Payload
	File      *os.File
	WaitGroup *sync.WaitGroup
	TarWriter *tar.Writer
}

// New create new writer client
func New(file *os.File) *Client {
	tarWriter := tar.NewWriter(file)

	return &Client{
		File:      file,
		Queue:     make(chan Payload),
		WaitGroup: &sync.WaitGroup{},
		TarWriter: tarWriter,
	}
}

// Start starting writing queue
func (client *Client) Start() {
	for payload := range client.Queue {
		client.Worker(payload)
	}
}

// Add function to add to writer payload
func (client *Client) Add(payload Payload) {
	client.WaitGroup.Add(1)
	client.Queue <- payload
}

// Worker writer process function
func (client *Client) Worker(payload Payload) error {
	defer client.WaitGroup.Done()

	// Set file tarball header
	header := new(tar.Header)
	header.Name = payload.Name
	header.Size = payload.Size
	header.Mode = 0766

	// write the header to the tarball archive
	if err := client.TarWriter.WriteHeader(header); err != nil {
		return err
	}

	// write conent to the tarball archive
	if _, err := io.Copy(client.TarWriter, payload.ReadCloser); err != nil {
		return err
	}

	return nil
}

// Close func to finish the queue
func (client *Client) Close() {
	client.WaitGroup.Wait()
	close(client.Queue)
	client.TarWriter.Close()
	client.File.Close()
}
