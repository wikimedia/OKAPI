package aws

import (
	"bytes"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// UploadFileToS3 uploading file to s3 bucket
func (connection *Connection) UploadFileToS3(path string, buffer *bytes.Buffer) (*s3.PutObjectOutput, error) {
	return connection.S3.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(connection.Options.Bucket),
		Key:         aws.String(path),
		Body:        bytes.NewReader(buffer.Bytes()),
		ContentType: aws.String(http.DetectContentType(buffer.Bytes())),
	})
}

// GetObjectFromS3 getting file bytes from s3 bucket
func (connection *Connection) GetObjectFromS3(path string) (*s3.GetObjectOutput, error) {
	return connection.S3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(connection.Options.Bucket),
		Key:    aws.String(path),
	})
}

// DeleteObjectFromS3 delete file from s3 bucket
func (connection Connection) DeleteObjectFromS3(path string) (*s3.DeleteObjectOutput, error) {
	return connection.S3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(connection.Options.Bucket),
		Key:    aws.String(path),
	})
}

// GetObjectRequestFromS3 getting file URL from s3 bucket
func (connection *Connection) GetObjectRequestFromS3(path string) (req *request.Request, output *s3.GetObjectOutput) {
	return connection.S3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(connection.Options.Bucket),
		Key:    aws.String(path),
	})
}

// ListS3BucketObjects listing objects in bucket
func (connection *Connection) ListS3BucketObjects() (*s3.ListObjectsV2Output, error) {
	return connection.S3.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(connection.Options.Bucket)})
}

// StreamFileToS3 stream the file to s3 bucket
func (connection *Connection) StreamFileToS3(path string, body io.Reader) (*s3manager.UploadOutput, error) {
	return connection.Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(connection.Options.Bucket),
		Key:    aws.String(path),
		Body:   body,
	})
}
